#自动导出restful api接口文档
官方文档(英文):http://apidocjs.com

官方文档(翻译):https://www.jianshu.com/p/9353d5cc1ef8


###1.安装nodejs

###2.安装apidoc
npm install -g apidoc

###3.安装apidoc-markdown
npm install -g apidoc-markdown

###4.自动导出restful api接口文档(html)
①cd $GOPATH/isoc/isoc_doc/auto_export_api_doc (cd 到工程目录下的doc/auto_export_api_doc文件夹)
②apidoc -i ./../../http/controller 

###5.根据生产的html接口文档生成markdown文档
apidoc-markdown -p ./doc -t markdown_template.md -o isoc_gateway.md

apidoc-markdown -p ./doc -t markdown_template.md -o all.md



