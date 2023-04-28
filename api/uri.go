// Copyright 2023 Planet Labs PBC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

const (
	crsEPSG3857 = "http://www.opengis.net/def/crs/EPSG/0/3857"
	wkssGoogle  = "http://www.opengis.net/def/wkss/OGC/1.0/GoogleMapsCompatible"
	tmsURIRoot  = "http://www.opengis.net/def/tilematrixset/OGC/1.0/"
)

func getTileMatrixSetURI(id string) string {
	return tmsURIRoot + id
}
