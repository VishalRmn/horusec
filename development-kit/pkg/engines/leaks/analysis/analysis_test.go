// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package analysis

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	engine "github.com/ZupIT/horusec-engine"
	"github.com/ZupIT/horusec/development-kit/pkg/cli_standard/config"
	"github.com/ZupIT/horusec/development-kit/pkg/enums/engine/advisories/leaks/regular"
	"github.com/stretchr/testify/assert"
)

func TestNewAnalysis(t *testing.T) {
	assert.IsType(t, NewAnalysis(config.NewConfig()), &Analysis{})
}

func TestAnalysis_StartAnalysis(t *testing.T) {
	t.Run("Should return success when read all example analysis and return 17 vulnerabilities", func(t *testing.T) {
		configs := config.NewConfig()
		configs.SetOutputFilePath("./leaks-tmp1.output.json")
		configs.SetProjectPath("../../examples")
		err := NewAnalysis(configs).StartAnalysis()
		assert.NoError(t, err)
		fileBytes, err := ioutil.ReadFile("./leaks-tmp1.output.json")
		data := []engine.Finding{}
		_ = json.Unmarshal(fileBytes, &data)
		assert.NoError(t, os.RemoveAll(configs.GetOutputFilePath()))
		assert.Equal(t, len(data), 19)
	})
	t.Run("Should return success when read analysis and return two vulnerabilities", func(t *testing.T) {
		configs := config.NewConfig()
		configs.SetOutputFilePath("./leaks-tmp2.output.json")
		configs.SetProjectPath("../../examples/leaks-hardcodedpass")
		err := NewAnalysis(configs).StartAnalysis()
		assert.NoError(t, err)
		fileBytes, err := ioutil.ReadFile("./leaks-tmp2.output.json")
		data := []engine.Finding{}
		_ = json.Unmarshal(fileBytes, &data)
		assert.NoError(t, os.RemoveAll(configs.GetOutputFilePath()))
		assert.Equal(t, len(data), 2)
	})
	t.Run("Should return success when read analysis and return empty vulnerabilities", func(t *testing.T) {
		configs := config.NewConfig()
		configs.SetOutputFilePath("./leaks-tmp3.output.json")
		err := NewAnalysis(configs).StartAnalysis()
		assert.NoError(t, err)
		fileBytes, err := ioutil.ReadFile("./leaks-tmp3.output.json")
		data := []engine.Finding{}
		_ = json.Unmarshal(fileBytes, &data)
		assert.NoError(t, os.RemoveAll(configs.GetOutputFilePath()))
		assert.Equal(t, len(data), 0)
	})
	t.Run("Should return error when create file", func(t *testing.T) {
		configs := config.NewConfig()
		configs.SetOutputFilePath("./////")
		err := NewAnalysis(configs).StartAnalysis()
		assert.Error(t, err)
	})
	t.Run("Should return error when get units in project path", func(t *testing.T) {
		configs := config.NewConfig()
		configs.SetOutputFilePath("./////")
		configs.SetProjectPath("./not exists path")
		err := NewAnalysis(configs).StartAnalysis()
		assert.Error(t, err)
	})
}

func TestAnalysis_StartRegularAnalysis(t *testing.T) {
	t.Run("Should return a vulnerability from PasswordExposedInHardcodeURL", func(t *testing.T) {
		configs := config.NewConfig()
		configs.SetOutputFilePath("./leaks-tmp4.output.json")
		configs.SetProjectPath("../../examples/leaks-hardcodedpass")
		err := NewAnalysis(configs).StartAnalysis()
		assert.NoError(t, err)
		fileBytes, err := ioutil.ReadFile("./leaks-tmp4.output.json")
		data := []engine.Finding{}
		_ = json.Unmarshal(fileBytes, &data)
		assert.NoError(t, os.RemoveAll(configs.GetOutputFilePath()))
		vulnCounter := 0
		for _, vuln := range data {
			if vuln.ID == regular.NewLeaksRegularPasswordExposedInHardcodedURL().ID {
				vulnCounter++
			}
		}
		assert.Equal(t, vulnCounter, 1)
	})

	t.Run("Should return a vulnerability from WPConfig", func(t *testing.T) {
		configs := config.NewConfig()
		configs.SetOutputFilePath("./leaks-tmp5.output.json")
		configs.SetProjectPath("../../examples/leaks-php")
		err := NewAnalysis(configs).StartAnalysis()
		assert.NoError(t, err)
		fileBytes, err := ioutil.ReadFile("./leaks-tmp5.output.json")
		data := []engine.Finding{}
		_ = json.Unmarshal(fileBytes, &data)
		assert.NoError(t, os.RemoveAll(configs.GetOutputFilePath()))
		vulnCounter := 0
		for _, vuln := range data {
			if vuln.ID == regular.NewLeaksRegularWPConfig().ID {
				vulnCounter++
			}
		}
		assert.Equal(t, 12, vulnCounter)
	})
}
