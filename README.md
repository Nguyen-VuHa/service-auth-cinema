# SERVICE AUTHENTICATION USING GOLANG

### Cấu trúc thư mục

```
- src
    - Constants
    - Controllers
    - DataLayers
      - Layer_1
      - Layer_2
      ...
    - DTO
    - Models
    - Public 
      - logs/...
    - Repositories
    - Routers
    - Services
       - Service_1
       - Service_2
       ...
    - Validations
    - ViewModels
    main.go
```

### Sử dụng cấu trúc MVVM và áp dụng các design pattern.

<p>
- Design Pattern: <br>
+ Repository Pattern <br>
+ Stragy Pattern 
</p>
 

### Tính năng:
[1. Đăng Nhập](#function-đăng-nhập) <br>
[2. Đăng ký](#function-đăng-ký) <br>
[3. Đăng nhập với Facebook]() <br>
[4. Đăng nhập với Google]() <br> 
[5. Đổi mật khẩu]() <br>
[6. Xác thực tài khoản]() <br>

## Function Đăng Nhập
* Workflow
```
B1. Client request -> server

B2. Server ghi nhận dữ liệu (Giải mã dữ liệu từ request) 
-> Kiểm tra dữ liệu (Kiểm tra các trường dữ liệu hợp lệ) (Router -> Controller)

B3. Xử lý dữ liệu từ B2
-> Kiểm tra Email đã đăng ký trong hệ thống (Services -> DataLayer)
-> compare password request user với password database (Compare password, ...) (Helpers)
-> Tạo Token cho user và lưu trữ lên Redis (Services -> Redis)

B4. Return response về client (Controller -> ViewModels)
```

## Function Đăng Ký
* Workflow
```
B1. Client request -> server

B2. Server ghi nhận dữ liệu (Giải mã dữ liệu từ request) 
-> Kiểm tra dữ liệu (Kiểm tra các trường dữ liệu hợp lệ) (Router -> Controller)

B3. Xử lý dữ liệu từ B2 
-> Kiểm tra Email đã tồn tại trong hệ thống (Services -> DataLayer)
-> Xử lý dữ liệu nhạy cảm của user (mã hoá Password, ...) (Helpers)
-> Ghi dữ liệu vào database (Services -> DataLayer)

B4. Return response về client (Controller -> ViewModels)
```
