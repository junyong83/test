// test
// test
// test
// test
// test
// test
// test
// haha

eroqgjpeqorkgeqrg

oauth2_proxy with Naver NSS 
===========================

// test
// test
// test

[oauth2_proxy](https://github.com/bitly/oauth2_proxy)는 다양한 OAuth provider을 사용하여 인증할 수 있게 한 reverse proxy 서버입니다.

이 프로젝트는 [NSS OAuth 2.0 API](http://wiki.navercorp.com/display/oapi/NSS+oAuth+2.0)를 추가하여 사내 서비스에 사번 인증을 가능하게 한 fork 버전으로, 변경 사항은 `naver` 브랜치에 적용되어 있습니다.

oauth2_proxy 대한 것은 다음 문서를 참고하세요: [oauth2_proxy_readme](README-oauth2-proxy.md)



### 사용방법
#### 선행 조건
  - go build 가 가능한 환경이 구축되어야 합니다. (GOROOT, GOPATH 설정되어 있어야 함)
  - NSS OAuth 2.0 API 를 사용할 수 있어야 합니다. 만약 ACL이 없다면 dev 환경으로 테스트는 가능합니다: [NSS API 사용신청](http://wiki.navercorp.com/display/oapi/API)

#### 실행하기
1. 소스코드 다운로드 및 빌드

    ```
    $ git clone https://oss.navercorp.com/NAVER-SEARCH/oauth2_proxy.git
    $ cd oauth2_proxy
    $ git checkout naver      # naver 브랜치에 적용되어 있습니다.
    $ go get -t
    $ go build .              # 빌드 파일 이름은 `oauth2_proxy`
    ```
2. proxy 서버 실행

    ```
    # 사용 예시 - proxy 서버 단독으로 80 포트를 가지고 띄울 때
    $ ./oauth2_proxy \  
            -client-id=MYSERVICE \              # NSS ACL 신청시 서비스 ID 
            -provider=naver \                   # NSS 사용을 위한 provider 이름
                                                # 이 값은 바꾸면 안됩니다.
            -upstream=http://127.0.0.1:8080 \   # proxy 요청을 보낼 하단 서버
            -redirect-url=http://csw000.nhnsystem.com:80 \ # OAuth redirect 주소
            -http-address=http://csw000.nhnsystem.com:80 \ # proxy 서버 주소 
            -skip-provider-button=true \        # provider에서 auth button 페이지 무시할때
                                                # 이 값은 바꾸면 안됩니다.
            -email-domain="navercorp.com" \     # 사번 인증시 유효한 사용자 email domain
            -cookie-secret="*" \                # 실행시 필요한 값. 변경 금지
            -client-secret="*" \                # 실행시 필요한 값. 변경 금지
            -cookie-secure=false                # 실행시 필요한 값. 변경 금지
    ```

####  앞단 웹서버 + oauth2 proxy 조합으로 실행하기
* oauth2_proxy 자체만으로도 서버 실행이 가능하나, 기존에 front 담당 웹서버가 있거나 SSL 등의 작업이 필요한 경우는
oauth2_proxy 앞에 웹서버를 놓고 proxy 기능만 수행하도록 합니다.
    - 아무런 설정값이 없을때 oauth2_proxy의 기본 포트는 4180
    - systemd 로 띄우려면 service config 파일을 참고하세요: [oauth2_proxy.service.example](contrib/oauth2_proxy.service.example)
    - cfg 파일로 이용할 수도 있습니다: [oauth2_proxy.cfg.example](contrib/oauth2_proxy.cfg.example)
* nginx 설정 예시

    ```
    # nginx.conf
    server {
        listen 80;
        server_name ooo.navercorp.com;

        location / {
            proxy_pass http://127.0.0.1:4180;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Scheme $scheme;
            proxy_connect_timeout 1;
            proxy_send_timeout 30;
            proxy_read_timeout 30;
        }
    }
    ```
* oauth2_proxy 실행

    ```
    # 사용 예시 - 상단 nginx 서버 아래 proxy로 띄울 때
    # NSS ACL 없이 개발환경에서 테스트하는 상황에는 제일 아래 3줄 필수 추가되어야 함
    $ ./oauth2_proxy \  
            -client-id=MYSERVICE \              # NSS ACL 신청시 서비스 ID 
            -provider=naver \                   # NSS 사용을 위한 provider 이름
                                                # 이 값은 바꾸면 안됩니다.
            -upstream=http://127.0.0.1:8080 \   # proxy 요청을 보낼 하단 서버
            -redirect-url=http://ooo.navercorp.com \ # OAuth redirect 주소. 상단 nginx access URL로 넣어주면 됨
            -http-address=127.0.0.1:4180 \      # proxy 서버 주소 (localhost:4180 일때는 생략가능)
            -skip-provider-button=true \        # provider에서 auth button 페이지 무시할때
                                                # 이 값은 바꾸면 안됩니다.
            -email-domain="navercorp.com" \     # 사번 인증시 유효한 사용자 email domain
            -cookie-secret="*" \                # 실행시 필요한 값. 변경 금지
            -client-secret="*" \                # 실행시 필요한 값. 변경 금지
            -cookie-secure=false                # 실행시 필요한 값. 변경 금지
            -login-url="https://nss-dev.navercorp.com/nweauthorize"
            -redeem-url="https://nss-api-dev.navercorp.com:5001/api/Auth/token"
            -validate-url="https://nss-api-dev.navercorp.com:5001/api/Auth/tokenInfo"
    ```

### 사용 사례
* 통합 로그검색시스템 Logiss : http://prod.logiss.navercorp.com/


### 참고
* oauth2_proxy 실행 옵션의 자세한 정보는 [oauth2_proxy README - command-line-options](README-oauth2-proxy.md#command-line-options)를 참고하세요.
* SSL에서 naver provider를 테스트 해 보지는 않았습니다. (사내 서비스 대부분이 non-SSL 환경인 상황)
* 잘못된 부분 수정 및 개선사항 적용이 있을때는 언제든지 PR 부탁드립니다~!
* 사용 사례도 알려주시면 추가하겠습니다.









