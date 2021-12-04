TEST1:
	C:\curl-7.80.0-win64-mingw\bin\curl.exe -v -L -X GET "http://localhost:8080/weather/?location=murmansk"
	curl -v -L -X GET "http://localhost:8080/weather/?location=murmansk"
	// тест ошибки
	C:\curl-7.80.0-win64-mingw\bin\curl.exe -v -L -X GET "http://localhost:8080/weather/"