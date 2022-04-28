package module

import (
	"embed"
	"path"
	"strconv"

	"github.com/cloudquery/cq-provider-sdk/cqproto"
	"github.com/hashicorp/go-hclog"
)

// InfoReader is called when the user executes a module, to get provider supported metadata about the given module
type InfoReader func(logger hclog.Logger, module string, prefferedVersions []uint32) (res InfoResult, err error)

// InfoResult is what the provider returns from an InfoReader request.
type InfoResult struct {
	// Data contains all cqproto.ModuleInfo incarnations in the requested/preferred versions (if any)
	Data map[uint32]cqproto.ModuleInfo
	// AvailableVersions is the all available versions supported by the provider (if any)
	AvailableVersions []uint32
}

// EmbeddedReader returns an InfoReader handler given a "moduleData" filesystem.
// The fs should have all the required files for the modules in basedir, as one subdirectory per module ID.
// Each subdirectory (for the module ID) should contain one subdirectory per protocol version.
// Each protocol-version subdirectory can contain multiple files.
// Example: moduledata/drift/1/file.hcl (where "drift" is the module name and "1" is the protocol version)
func EmbeddedReader(moduleData embed.FS, basedir string) InfoReader {
	return func(logger hclog.Logger, module string, prefferedVersions []uint32) (InfoResult, error) {
		var (
			res InfoResult
			err error
		)

		res.AvailableVersions, err = availableVersions(moduleData, path.Join(basedir, module))
		if err != nil {
			return res, err
		}

		res.Data = make(map[uint32]cqproto.ModuleInfo, len(prefferedVersions))

		for _, v := range prefferedVersions {
			dir := path.Join(basedir, module, strconv.FormatInt(int64(v), 10)) // <basedir>/<module>/<version>/
			data, err := flatFiles(moduleData, dir, "")
			if err != nil {
				return res, err
			}
			if len(data) == 0 {
				continue
			}

			inf := cqproto.ModuleInfo{
				Files: data,
			}
			res.Data[v] = inf
		}

		if len(prefferedVersions) > 0 && len(res.Data) == 0 {
			logger.Warn("received unsupported module info request", "module", module, "preferred_versions", prefferedVersions)
		}

		return res, nil
	}
}

func availableVersions(moduleData embed.FS, dir string) ([]uint32, error) {
	versionDirs, err := moduleData.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	versions := make([]uint32, 0, len(versionDirs))
	for _, f := range versionDirs {
		if !f.IsDir() {
			continue
		}
		vInt, err := strconv.ParseUint(f.Name(), 10, 32)
		if err != nil {
			continue
		}
		versions = append(versions, uint32(vInt))
	}

	return versions, nil
}

func flatFiles(moduleData embed.FS, dir, prefix string) ([]*cqproto.ModuleFile, error) {
	files, err := moduleData.ReadDir(dir)
	if err != nil {
		return nil, nil
	}

	var ret []*cqproto.ModuleFile
	for _, f := range files {
		name := path.Join(dir, f.Name())

		if !f.IsDir() {
			data, err := moduleData.ReadFile(name)
			if err != nil {
				return nil, err
			}
			ret = append(ret, &cqproto.ModuleFile{
				Name:     path.Join(prefix, f.Name()),
				Contents: data,
			})
			continue
		}

		// recurse and read subdirs
		sub, err := flatFiles(moduleData, name, f.Name())
		if err != nil {
			return nil, err
		}
		ret = append(ret, sub...)
	}

	return ret, nil
}
