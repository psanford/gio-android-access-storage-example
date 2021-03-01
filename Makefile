PLATFORM_JAR=$(ANDROID_HOME)/platforms/android-30/android.jar

gio-android-access-storage-example: sanford_jni.jar $(wildcard *.go)
	go run gioui.org/cmd/gogio -target android .

sanford_jni.jar: Jni.java
	mkdir -p classes
	javac -cp $(PLATFORM_JAR) -sourcepath $(PLATFORM_JAR) -d classes $^
	jar cf $@ -C classes .
	rm -rf classes
