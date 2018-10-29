package main

import (
	"archive/zip"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/privatix/dapp-openvpn/statik"
)

/*

	Package structure:

	rootPath:
		productPath:
			bin:
				-installer
				-adapter
			template:
				-product.client.json or product.agent.json
				-template.offering.json
				-template.access.json
			data:
			config:
				-adapter.config.json

*/

const (
	id        = "4b26dc82-ffb6-4ff1-99d8-f0eaac0b0532"
	commit    = "commit"
	release   = "release"
	pkg       = "pkg"
	bin       = "bin"
	tmp       = "tmp"
	build     = "build"
	tags      = "-tags=notest"
	flags     = "-ldflags"
	exeSuffix = ".exe"
	zipSuffix = ".zip"

	templatesSrc  = "/package/template"
	configSrc     = "/package/config"
	clientProduct = "/product/product.client.json"
	agentProduct  = "/product/product.agent.json"

	// Permissions.
	pathPerm = 0755
	filePerm = 0644
)

var (
	adapterBin   = "adapter"
	installerBin = "installer"
	adapterPkg   = adapterBin
	installerPkg = installerBin

	repoPath string
	binPath  string
	ldFlags  string
	target   string

	// "xgo" targets
	targets = map[string]string{
		"linux":   "--targets=linux/amd64",
		"windows": "--targets=windows/amd64",
		"macos":   "--targets=darwin/amd64",
	}

	// If is true, then use the "xgo" to create binary files.
	xgo bool

	// If is true, then the package will be created for a agent.
	agent bool

	commands = map[string]*command{
		// Gets last commit.
		commit: {"git", []string{"rev-list", "-1", "HEAD"}, ""},
		// Gets version tag.
		release: {"git", []string{"tag", "-l",
			"--points-at", "HEAD"}, ""},
		// Gets repository pkg.
		pkg: {"go", []string{"list"}, ""},
	}

	versionPackage string
	minCoreVersion string
	maxCoreVersion string

	zipName string
)

type command struct {
	app    string
	args   []string
	result string
}

type descriptor struct {
	Name              string `json:"name"`
	ID                string `json:"id"`
	Version           string `json:"version"`
	MinCoreAppVersion string `json:"min_core_app_version"`
	MaxCoreAppVersion string `json:"max_core_app_version"`
	Signature         string `json:"signature"`
	Hash              string `json:"hash"`
}

func main() {
	var err error

	flag.StringVar(&versionPackage, "version", "undefined",
		"Product package distributive version.")
	flag.StringVar(&minCoreVersion, "min_core_version", "undefined",
		"Minimal version of Privatix core application for"+
			" compatibility.")
	flag.StringVar(&maxCoreVersion, "max_core_version", "",
		"Maximum version of Privatix core application for"+
			" compatibility.")
	flag.StringVar(&target, "os", "",
		"Target OS: linux, windows or macos (xgo usage). "+
			"If is empty, a package will be created "+
			"for the current operating system.")

	flag.BoolVar(&agent, "agent", false, "Whether to install agent.")

	jsonBlob := flag.String(
		"keystore", "", "Full path to JSON private key file.")

	auth := flag.String("auth", "",
		"Password to decrypt JSON private key.")

	flag.Parse()

	checkVersionFlags()

	pk := privateKeyFromFile(*jsonBlob, *auth)

	checkTargetOS()

	// Get repository full path.
	repoPath, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	binPath = filepath.Join(repoPath, build, tmp, id, bin)

	if err := os.RemoveAll(filepath.Join(repoPath, build)); err != nil {
		panic(err)
	}

	// Erase temporary directory.
	defer os.RemoveAll(filepath.Join(repoPath, build, tmp))

	// Need to create a bin folder in advance so that xgo does not do it.
	// A "xgo" creates a folder, but its owner is a root user.
	// It is impossible to remove it later.
	// {repo}/build/tmp/product/{id}/bin
	if err := os.MkdirAll(binPath, pathPerm); err != nil {
		panic(err)
	}

	getParams()
	runCommands()
	copyFiles()

	archive := filepath.Join(repoPath, build, zipName+zipSuffix)
	dc := filepath.Join(repoPath, build, zipName+".json")

	compress(filepath.Join(repoPath, build, tmp, id), archive)
	hash, signature := sign(archive, pk)

	makeDescriptor(signature, hash, dc)
}

func checkVersionFlags() {
	if versionPackage == "" {
		panic("product package distributive version must not be empty")
	}

	if minCoreVersion == "" {
		panic("minimal version of Privatix core application for" +
			" compatibility must not be empty")
	}
}

func privateKeyFromFile(file, pass string) *ecdsa.PrivateKey {
	blob, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	key, err := keystore.DecryptKey(blob, pass)
	if err != nil {
		panic(err)
	}
	return key.PrivateKey
}

func sign(file string, pk *ecdsa.PrivateKey) (h string, s []byte) {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	hash := crypto.Keccak256(raw)
	sig, err := crypto.Sign(hash, pk)
	if err != nil {
		panic(err)
	}

	return common.BytesToHash(hash).String(), sig
}

func makeDescriptor(signature []byte, hash, output string) {
	dc := &descriptor{
		Name:              "openvpn",
		ID:                id,
		Version:           versionPackage,
		MinCoreAppVersion: minCoreVersion,
		MaxCoreAppVersion: maxCoreVersion,
		Signature:         base64.URLEncoding.EncodeToString(signature),
		Hash:              hash,
	}

	raw, _ := json.Marshal(dc)
	if err := ioutil.WriteFile(output, raw, filePerm); err != nil {
		panic(err)
	}
}

func checkXgo() bool {
	path, err := exec.LookPath("xgo")
	if err != nil {
		fmt.Printf("didn't find 'xgo' executable\n")
		return false
	}
	fmt.Printf("'xgo' executable is in '%s'\n", path)
	return true
}

func checkTargetOS() {
	if target != "" {
		if _, ok := targets[target]; ok {
			// check "xgo" executable
			if !checkXgo() {
				panic("xgo must be executable")
			}
			xgo = true

			zipName = target[:3]
		} else {
			panic("unsupported operation system")
		}
	} else {
		switch runtime.GOOS {
		case "windows", "linux":
			zipName = runtime.GOOS[:3]
		case "darwin":
			zipName = "mac"
		default:
			panic("unsupported operation system")
		}
	}

	zipName = zipName + "_" + versionPackage
}

func getParams() {
	if runtime.GOOS == "windows" {
		adapterPkg += exeSuffix
		installerPkg += exeSuffix
	}

	for k, v := range commands {
		val, err := exec.Command(v.app, v.args...).Output()
		if err != nil {
			panic(err)
		}
		commands[k].result = strings.TrimSpace(string(val))
	}

	ldFlags = fmt.Sprintf(`-X main.Commit=%s -X main.Version=%s`,
		commands[commit].result, commands[release].result)
}

// buildFlags: build -ldflags "-X main.Commit={commit} -X main.Version={release}"
// -tags=notest -o {full path to an output file} {go pkg to build}
func buildFlags(binName, packageName string) []string {
	return []string{"build", flags, ldFlags, tags, "-o",
		filepath.Join(binPath, binName),
		filepath.Join(commands[pkg].result, packageName)}
}

// buildXGoFlags: -ldflags "-X main.Commit={commit} -X main.Version={release}"
// -tags=notest -out {output file name} --targets={target os and architecture}
// {full path to pkg to build}
func buildXGoFlags(binName, packageName string) []string {
	return []string{flags, ldFlags, tags, "-out",
		filepath.Join(build, tmp, id, bin, binName),
		targets[target], filepath.Join(repoPath, packageName)}
}

func runCommands() {
	adapterBuild := &command{"go",
		buildFlags(adapterBin, adapterPkg), ""}
	installerBuild := &command{"go",
		buildFlags(installerBin, installerPkg), ""}

	if xgo {
		adapterBuild = &command{"xgo",
			buildXGoFlags(adapterBin, adapterPkg), ""}
		installerBuild = &command{"xgo",
			buildXGoFlags(installerBin, installerPkg), ""}
	}

	run := []*command{
		{"dep", []string{"ensure"}, ""},
		{"go", []string{"generate",
			"." + string(filepath.Separator) + "..."}, ""},
		adapterBuild,
		installerBuild,
	}

	for _, v := range run {
		if err := exec.Command(v.app, v.args...).Run(); err != nil {
			panic(err)
		}
	}

	if xgo {
		renameXGOFiles()
	}
}

// renameXGOFiles renames binaries after xgo building.
func renameXGOFiles() {
	files, err := ioutil.ReadDir(filepath.Join(
		repoPath, build, tmp, id, bin))
	if err != nil {
		panic(err)
	}

	for _, v := range files {
		if v.IsDir() {
			continue
		}

		var name string

		if strings.HasPrefix(v.Name(), adapterPkg) {
			name = adapterPkg
		} else if strings.HasPrefix(v.Name(), installerPkg) {
			name = installerPkg
		} else {
			continue
		}

		// If xgo and windows os then adds ".exe" suffix to binaries.
		if target == "windows" {
			name += exeSuffix
		}

		if err := os.Rename(
			filepath.Join(repoPath, build, tmp, id, bin, v.Name()),
			filepath.Join(repoPath, build, tmp,
				id, bin, name)); err != nil {
			panic(err)
		}
	}
}

// copyFiles copies files from statik filesystem to package temp folder.
func copyFiles() {
	rootPackagePath := filepath.Join(repoPath, build, tmp, id)

	templatesDst := filepath.Join(rootPackagePath, "template")
	configDst := filepath.Join(rootPackagePath, "config")

	files := map[string]string{
		templatesSrc: templatesDst,
		configSrc:    configDst,
	}

	if agent {
		files[agentProduct] = filepath.Join(
			templatesDst, "product.agent.json")
	} else {
		files[clientProduct] = filepath.Join(
			templatesDst, "product.client.json")
	}

	for src, dst := range files {
		file, err := statik.OpenFile(src)
		if err != nil {
			panic(err)
		}

		fileInfo, err := file.Stat()
		if err != nil {
			panic(err)
		}

		copy(src, dst, fileInfo)
	}
}

// compress compresses package.
func compress(source, target string) {
	zipfile, err := os.Create(target)
	if err != nil {
		panic(err.Error())
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		panic(err.Error())
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	if err := filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(
					baseDir, strings.TrimPrefix(path, source))
			}

			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}

			writer, err := archive.CreateHeader(header)
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			return err
		}); err != nil {
		panic(err.Error())
	}
}

func copy(src, dst string, info os.FileInfo) error {
	if info.IsDir() {
		return dirCopy(src, dst)
	}
	return fileCopy(src, dst)
}

func fileCopy(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), pathPerm); err != nil {
		return err
	}

	f, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	s, err := statik.OpenFile(src)
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = io.Copy(f, s)
	return err
}

func dirCopy(srcDir, dstDir string) error {
	if err := os.MkdirAll(dstDir, pathPerm); err != nil {
		return err
	}

	contents, err := statik.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, content := range contents {
		var base string
		if content.IsDir() {
			base = filepath.Base(content.Name())

		} else {
			base = content.Name()
		}

		cs := filepath.Join(srcDir, base)
		cd := filepath.Join(dstDir, base)

		if err := copy(cs, cd, content); err != nil {
			return err
		}
	}
	return nil
}
