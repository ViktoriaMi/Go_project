TEST1:
	#C:\curl-7.80.0-win64-mingw\bin\curl.exe -v -L -X GET "http://localhost:8080/weather/?location=murmansk"
	curl -v -L -X GET "http://localhost:8080/weather/?location=murmansk"
	// тест ошибки
	C:\curl-7.80.0-win64-mingw\bin\curl.exe -v -L -X GET "http://localhost:8080/weather/"
	export APIKEY="d903717989b333890700e2644c0c7a8e"
	$(shell sourse .env)
	$(shell sourse .env) &&\
	sourse .env

TEST2:
	$(shell sourse .env)
	curl -v -L -X GET "http://localhost:8080/weather/?location=murmansk"