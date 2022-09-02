## 安装opencv



## CentOS 7 安装 Cmake3

```
yum -y install cmake3
mv /usr/bin/cmake /usr/bin/cmake2
ln -s /usr/bin/cmake3 /usr/bin/cmake
```



## CentOS 7 安装 gcc 8.3.0

```
curl -LJO http://mirror.linux-ia64.org/gnu/gcc/releases/gcc-8.3.0/gcc-8.3.0.tar.gz
tar xvfz gcc-8.3.0.tar.gz && cd gcc-8.3.0
./contrib/download_prerequisites
./configure --disable-multilib --enable-languages=c,c++ --prefix=/usr/local
make -j5
make -j install
```



### 编译安装 opencv

只需编译这个3个库 core, calib3d, imgproc

```
tar xvfz opencv-4.5.5.tar.gz
cd opencv-4.5.5
mkdir -p build && cd build

cmake -D CMAKE_BUILD_TYPE=RELEASE -D CMAKE_INSTALL_PREFIX=/usr/local -D BUILD_DOCS=OFF -D BUILD_EXAMPLES=OFF -D BUILD_TESTS=OFF -D BUILD_PERF_TESTS=OFF -D BUILD_opencv_java=NO -D BUILD_opencv_python=NO -D BUILD_opencv_python2=NO -D BUILD_opencv_python3=NO -D WITH_JASPER=OFF -D WITH_TBB=ON -D BUILD_SHARED_LIBS=ON -DOPENCV_GENERATE_PKGCONFIG=ON -DWITH_IPP=OFF -DBUILD_LIST=core,calib3d,imgproc ..
cmake --build .

sudo make install
```
