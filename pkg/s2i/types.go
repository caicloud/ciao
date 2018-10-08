// Copyright 2018 Caicloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package s2i

import (
	"github.com/caicloud/ciao/pkg/types"
)

// Interface is the interface for s2i.
type Interface interface {
	SourceToImage(code string, parameter *types.Parameter) (string, error)
	Cleanup(parameter *types.Parameter) (error)
}
