/*
Copyright © 2021 Bold City Software

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

env.go - this file contains logic for dealing with environments and variables
*/
package napenv

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/davesheldon/nap/naputil"
	"gopkg.in/yaml.v2"
)

func AddEnvironmentFromPath(workingDirectory string, environmentFileName string, existing map[string]string) (map[string]string, error) {
	if path.Ext(environmentFileName) != ".yml" && path.Ext(environmentFileName) != ".yaml" {
		environmentFileName = environmentFileName + ".yml"
	}

	originalFileName := environmentFileName

	if len(originalFileName) > 0 {
		if exists, _ := naputil.FileExists(environmentFileName); !exists {
			// try and find it relative to the target path
			environmentFileName = filepath.Join(workingDirectory, "..", "env", originalFileName)
		}

		if exists, _ := naputil.FileExists(environmentFileName); !exists {
			// try and find it relative to the target path
			environmentFileName = filepath.Join(workingDirectory, "env", originalFileName)
		}

		if exists, _ := naputil.FileExists(environmentFileName); !exists {
			// try and find it relative to the target path
			environmentFileName = filepath.Join(workingDirectory, originalFileName)
		}

		if exists, err := naputil.FileExists(environmentFileName); !exists {
			return existing, fmt.Errorf("environment '%s' not found.", originalFileName)
		} else if err != nil {
			return existing, fmt.Errorf("cannot read '%s'. %e", originalFileName, err)
		}

		configData, err := os.ReadFile(environmentFileName)
		if err != nil {
			return existing, fmt.Errorf("cannot open '%s'. %e", originalFileName, err)
		}

		subMap := make(map[string]string)

		err = yaml.Unmarshal(configData, &subMap)
		if err != nil {
			return existing, fmt.Errorf("cannot parse '%s'. %e", originalFileName, err)
		}

		for k, v := range subMap {
			existing[k] = v
		}
	}

	return existing, nil
}
