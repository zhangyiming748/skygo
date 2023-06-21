# <%= project.name %> v<%= project.version %>


<% if (prepend) { -%>
<%- prepend %>
<% } -%>
<% Object.keys(data).forEach(function (group) { -%>
# <%= group %>

<% Object.keys(data[group]).forEach(function (sub) { -%>
# <%-: data[group][sub][0].type | upcase %> <%= data[group][sub][0].url %>

<%-: data[group][sub][0].description | undef %>

<% if (data[group][sub][0].header && data[group][sub][0].header.fields.Header.length) { -%>
请求头参数列表:

| 名称    | 类型      |必填       | 说明                          |
|---------|-----------|-----------|--------------------------------------|
<% data[group][sub][0].header.fields.Header.forEach(function (header) { -%>
| <%- header.field %>			| <%- header.type %>			| <%- header.optional ? 'no' : 'yes' %>  | <%- header.description %>							|
<% }); //forech parameter -%>
<% } //if parameters -%>
<% if (data[group][sub][0].parameter && data[group][sub][0].parameter.fields.Parameter.length) { -%>

+ Parameters
<% data[group][sub][0].parameter.fields.Parameter.forEach(function (param) { -%>
    + <%- param.field %> (<%- param.type %>,<%- param.optional ? 'optional' : '' %>) ... <%- param.description %><%if(param.defaultValue){%><br/>默认值:`<%- param.defaultValue %>`<%}%>  <%if(param.allowedValues){%> <br/>可选值: <% Object.keys(param.allowedValues).forEach(function (subAllowValue) {%>`<%-subAllowValue %>`,<%})}%>
<%});//forech parameter-%>
<% } //if parameters-%>
<%/*请求参数示例*/%>
<% if (data[group][sub][0].parameter && data[group][sub][0].parameter.examples) { -%>
+ Request
<% data[group][sub][0].parameter.examples.forEach(function (example) { -%>
<%- example.content %>
<% }); //foreach example -%>
<% } //if example -%>
<%/*请求成功示例*/%>
<% if (data[group][sub][0].success && data[group][sub][0].success.examples && data[group][sub][0].success.examples.length) { -%>
+ Response 200
<% data[group][sub][0].success.examples.forEach(function (example) { -%>
<%- example.content %>
<% }); //foreach success example -%>
<% } //if examples -%>
<% if (data[group][sub][0].error && data[group][sub][0].error.examples && data[group][sub][0].error.examples.length) { -%>

<% data[group][sub][0].error.examples.forEach(function (example) { -%>
<%= example.title %>

```<%= example.type %>
<%- example.content %>
```
<% }); //foreach error example -%>
<% } //if examples -%>
<% }); //foreach sub  -%>
<% }); //foreach group -%>

