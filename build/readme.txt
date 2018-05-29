
将 src 目录中的所有子账单文件，转换为 csv 格式输出到 dst 目录

src： 默认为当前目录中的 src 文件夹
dst： 默认为当前目录中的 dst 文件夹


命令行执行；
	billconverter -h
		Usage of billconverter:
  			-dst string
        			dst folder (default "./dst")
  			-src string
        			src folder (default "./src")
	
	billconverter	
		/*src, dst 默认在当前目录中*/

		
	billconverter -src="你/的/账/单/目/录" -dst="输/出/目/录"




**** 注：命令行执行结尾 如看到： 
	
	INFO: all file convert successed. 代表所有文件转换成功

	ERROR: ...... 61188801.txt ......  代表 这个 61188801.txt 账单转换失败，
		请截屏，并将此账单文件传给开发人员。