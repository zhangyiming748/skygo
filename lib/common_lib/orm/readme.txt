PhalconStyle是一种url参数请求风格，按照预定义的格式，前端拼接好请求参数，后端可以采用一致的解析器解析。

The request URL's params parts describe the conditions we need focus on.

Params [where] : it must be a standard json string describe the find conditions, it should use with operator like "e"、"gt" egs.
				for example:  ?where={"sn":{"e":"89860617070021807390"}}
Params [fields] : it must be a  string separated  by comma, describe the fields that must be shown.
				for example: ?where={"sn":{"e":"89860617070021807390"}}&fields=sn,brand
Params [sort]:  it must be a standard json string describe the order of the result. -1 represents descend, 1 represents Ascend
				for example: ?where={"brand":%20"{"e":"BYD"}}&limit=10&sort={"id":-1}
Params [offset] :  it is a number that tell us how many records should be skipped before search.
				for example: ?where={"brand":"{"e":"BYD"}}&offset=2&limit=10
Params [limit] :  it is a number that tell us how many records should be return
				for example: ?where={"brand":"{"e":"BYD"}}&offset=2&limit=10
Params [include] : it must be a string separated by comma, describe the corresponding models that must be search together
				for example: ?where={"brand":"{"e":"BYD"}}&offset=2&limit=10&include=vehicle_location