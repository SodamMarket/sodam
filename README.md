# sodam

Nortion 프로젝트 페이지
https://www.notion.so/teamtam/09ce4b7b7a704a50ad2be09381810cff

<h1>1. DB</h1>
<h2>cockroachDB</h2>
https://www.cockroachlabs.com/docs/stable/install-cockroachdb-windows.html

다운 받은 cockroach 파일을 $GOROOT or $GOPATH에 넣고 아래의 명령어를 입력해주세요.
추천경로 \$GOROOT\bin

<pre><code>cockroach start --insecure --host=localhost</pre></code>

<h1>프로젝트 다운</h1>
<pre><code>go get -u github.com/SodamMarket/sodam</pre></code>

<h2>패키지 다운</h2>
<pre><code>go mod download</pre></code>

<h2>실행</h2>
<pre><code>go run main.go</pre></code>

추가적으로 해당 코드는 vscode에서 작업하였으며 Rest를 테스트하기 위해 Rest Client라는 패키지를 설치하였습니다.
