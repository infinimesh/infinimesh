/*
Copyright © 2021-2023 Infinite Devices GmbH, Nikita Ivanovski info@slnt-opp.xyz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package infinimesh

type ContextKey string

const INFINIMESH_ROOT_CLAIM = "root"
const INFINIMESH_ACCOUNT_CLAIM = "account"

const INFINIMESH_DEVICES_CLAIM = "devices"
const INFINIMESH_POST_STATE_ALLOWED_CLAIM = "post"

const InfinimeshRootCtxKey = ContextKey(INFINIMESH_ROOT_CLAIM)
const InfinimeshAccountCtxKey = ContextKey(INFINIMESH_ACCOUNT_CLAIM)
const InfinimeshDevicesCtxKey = ContextKey(INFINIMESH_DEVICES_CLAIM)
const InfinimeshPostAllowedCtxKey = ContextKey(INFINIMESH_POST_STATE_ALLOWED_CLAIM)
