// Copyright © 2023 OpenIM open source community. All rights reserved.
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

package common

var Msg = map[int]string{
	SUCCESS:        "success",
	ERROR:          "error",
	INVALID_PARAMS: "request params error",
}

func GetMsg(code int) string {
	msg, ok := Msg[code]
	if ok {
		return msg
	}

	return Msg[ERROR]
}
