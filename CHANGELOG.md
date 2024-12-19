# Changelog

## [v1.6.4](https://github.com/fabiolb/fabio/tree/v1.6.4) (2024-11-27)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.6.3...v1.6.4)

**Closed issues:**

- CI pipeline to run testsuite [\#949](https://github.com/fabiolb/fabio/issues/949)
- Fabio not exporting all the metrics with prometheus [\#947](https://github.com/fabiolb/fabio/issues/947)
- certificates - cert and ca chain/intermediate [\#946](https://github.com/fabiolb/fabio/issues/946)
- This repository is unmaintaned. [\#944](https://github.com/fabiolb/fabio/issues/944)
- CVE-2023-44487 HTTP/2 rapid reset [\#939](https://github.com/fabiolb/fabio/issues/939)
- TCP no route - cant balance tcp [\#936](https://github.com/fabiolb/fabio/issues/936)
- windows: setting logging path in fabio properties [\#920](https://github.com/fabiolb/fabio/issues/920)
- Port range in the proxy.addr [\#529](https://github.com/fabiolb/fabio/issues/529)

**Merged pull requests:**

- Add GoReleaser workflow. [\#958](https://github.com/fabiolb/fabio/pull/958) ([tristanmorgan](https://github.com/tristanmorgan))
- go-kit/kit/log go-kit/log [\#956](https://github.com/fabiolb/fabio/pull/956) ([tristanmorgan](https://github.com/tristanmorgan))
- Update go-retryablehttp to fix warning. [\#954](https://github.com/fabiolb/fabio/pull/954) ([tristanmorgan](https://github.com/tristanmorgan))
- Update and try fix GH Pages publish action. [\#953](https://github.com/fabiolb/fabio/pull/953) ([tristanmorgan](https://github.com/tristanmorgan))
- Update go version, test binaries and package versions. [\#952](https://github.com/fabiolb/fabio/pull/952) ([tristanmorgan](https://github.com/tristanmorgan))
- Remove vendored modules in favour of go mod. [\#951](https://github.com/fabiolb/fabio/pull/951) ([tristanmorgan](https://github.com/tristanmorgan))
- Update Github runner image. [\#950](https://github.com/fabiolb/fabio/pull/950) ([tristanmorgan](https://github.com/tristanmorgan))
- fix: close resp body [\#945](https://github.com/fabiolb/fabio/pull/945) ([testwill](https://github.com/testwill))
- Trim leading and trailing spaces from service tags [\#943](https://github.com/fabiolb/fabio/pull/943) ([logocomune](https://github.com/logocomune))
- Fix doubled download a Vault file [\#942](https://github.com/fabiolb/fabio/pull/942) ([aleksraiden](https://github.com/aleksraiden))
- Remove deprecated ioutil [\#940](https://github.com/fabiolb/fabio/pull/940) ([tristanmorgan](https://github.com/tristanmorgan))
- Dockerfile: add CAP\_NET\_BIND\_SERVICE+eip to fabio to allow running as root [\#938](https://github.com/fabiolb/fabio/pull/938) ([Kamilcuk](https://github.com/Kamilcuk))
- Consul registry performance improvements [\#928](https://github.com/fabiolb/fabio/pull/928) ([ddreier](https://github.com/ddreier))
- \[Docs\] Fix wrong parameter name [\#914](https://github.com/fabiolb/fabio/pull/914) ([KEANO89](https://github.com/KEANO89))
- Updating grpc handler to gracefully close backend connections [\#913](https://github.com/fabiolb/fabio/pull/913) ([nathanejohnson](https://github.com/nathanejohnson))

## [v1.6.3](https://github.com/fabiolb/fabio/tree/v1.6.3) (2022-12-09)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.6.2...v1.6.3)

**Implemented enhancements:**

- Feature request: Make source links in ui interface clickable [\#901](https://github.com/fabiolb/fabio/issues/901)

**Closed issues:**

- Ignore host=dst when backend is https [\#916](https://github.com/fabiolb/fabio/issues/916)
- poll new feature requests [\#910](https://github.com/fabiolb/fabio/issues/910)
- Fabio Clustering. [\#668](https://github.com/fabiolb/fabio/issues/668)

**Merged pull requests:**

- Disable BGP functionality on Windows since gobgp does not support this. [\#919](https://github.com/fabiolb/fabio/pull/919) ([nathanejohnson](https://github.com/nathanejohnson))
- updating CHANGELOG [\#918](https://github.com/fabiolb/fabio/pull/918) ([nathanejohnson](https://github.com/nathanejohnson))
- Don't use "dst" literal as sni name when host=dst is specified on https backends [\#917](https://github.com/fabiolb/fabio/pull/917) ([nathanejohnson](https://github.com/nathanejohnson))
- add feature to advertise anycast addresses via BGP [\#909](https://github.com/fabiolb/fabio/pull/909) ([nathanejohnson](https://github.com/nathanejohnson))
- Change the shutdown procedure to deregister fabio from the registry and then shutdown the proxy [\#908](https://github.com/fabiolb/fabio/pull/908) ([martinivanov](https://github.com/martinivanov))
- Feature/source link [\#907](https://github.com/fabiolb/fabio/pull/907) ([KTruesdellENA](https://github.com/KTruesdellENA))

## [v1.6.2](https://github.com/fabiolb/fabio/tree/v1.6.2) (2022-09-13)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.6.1...v1.6.2)

**Closed issues:**

- Update TLS cipher parser to include modern ciphers [\#903](https://github.com/fabiolb/fabio/issues/903)
- Custom behavior for the situation when the service has no healthy instances [\#898](https://github.com/fabiolb/fabio/issues/898)

**Merged pull requests:**

- update README for v1.6.2 release [\#905](https://github.com/fabiolb/fabio/pull/905) ([nathanejohnson](https://github.com/nathanejohnson))
- Updating TLS cipher config parser to include TLS 1.3 constants. [\#904](https://github.com/fabiolb/fabio/pull/904) ([nathanejohnson](https://github.com/nathanejohnson))

## [v1.6.1](https://github.com/fabiolb/fabio/tree/v1.6.1) (2022-07-19)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.6.0...v1.6.1)

**Implemented enhancements:**

- Multi-DC fabio [\#115](https://github.com/fabiolb/fabio/issues/115)

**Fixed bugs:**

- Crash: invalid log msg: http2: panic serving CLIENT\_IP:CLIENT\_PORT: runtime error: index out of range \[-1\] [\#872](https://github.com/fabiolb/fabio/issues/872)

**Closed issues:**

- admin UI Overrides not working [\#886](https://github.com/fabiolb/fabio/issues/886)
- Panic on created prometheus metric name [\#878](https://github.com/fabiolb/fabio/issues/878)
- Crash on route update: panic: runtime error: index out of range, diffmatchpatch.\(\*DiffMatchPatch\).DiffCharsToLines [\#873](https://github.com/fabiolb/fabio/issues/873)
- Experiencing 502's [\#862](https://github.com/fabiolb/fabio/issues/862)
- Fabio immediately drop routes when consul agent unavailable for a while [\#861](https://github.com/fabiolb/fabio/issues/861)
- \[proxy/tls\] Update supported TLS versions and cipher suites [\#858](https://github.com/fabiolb/fabio/issues/858)
- JSON schema is incorrect in website Dest should be Dst [\#852](https://github.com/fabiolb/fabio/issues/852)
- \[question\] URL for TLS destination [\#850](https://github.com/fabiolb/fabio/issues/850)
- \[Feature\] Possibility of adding arm/arm64 docker builds. [\#833](https://github.com/fabiolb/fabio/issues/833)

**Merged pull requests:**

- Release/v1.6.1 [\#897](https://github.com/fabiolb/fabio/pull/897) ([nathanejohnson](https://github.com/nathanejohnson))
- setting sni to match host [\#896](https://github.com/fabiolb/fabio/pull/896) ([KTruesdellENA](https://github.com/KTruesdellENA))
- Update random picker to use math/rand's Intn function [\#893](https://github.com/fabiolb/fabio/pull/893) ([nathanejohnson](https://github.com/nathanejohnson))
- add configurable grpc message sizes to \#632 [\#890](https://github.com/fabiolb/fabio/pull/890) ([nathanejohnson](https://github.com/nathanejohnson))
- add tls13 [\#889](https://github.com/fabiolb/fabio/pull/889) ([nathanejohnson](https://github.com/nathanejohnson))
- update materialize bits. see issue \#886 [\#888](https://github.com/fabiolb/fabio/pull/888) ([nathanejohnson](https://github.com/nathanejohnson))
- Moved admin UI assets to use go embed [\#885](https://github.com/fabiolb/fabio/pull/885) ([nathanejohnson](https://github.com/nathanejohnson))
- update the custom css [\#884](https://github.com/fabiolb/fabio/pull/884) ([KTruesdellENA](https://github.com/KTruesdellENA))
- Bump go-diff dependency version to 1.2.0.  Fixes \#873 [\#881](https://github.com/fabiolb/fabio/pull/881) ([nathanejohnson](https://github.com/nathanejohnson))
- bump HUGO version to 0.101.0 [\#880](https://github.com/fabiolb/fabio/pull/880) ([nathanejohnson](https://github.com/nathanejohnson))
- add docs from PR \#854 to fabio.properties [\#879](https://github.com/fabiolb/fabio/pull/879) ([nathanejohnson](https://github.com/nathanejohnson))
- Build multi-arch Docker images for amd64 and arm64 architectures [\#875](https://github.com/fabiolb/fabio/pull/875) ([vamc19](https://github.com/vamc19))
- Fix x-forwarded-for header processing for ws connections [\#860](https://github.com/fabiolb/fabio/pull/860) ([bn0ir](https://github.com/bn0ir))
- Update registry.backend.md [\#854](https://github.com/fabiolb/fabio/pull/854) ([webmutation](https://github.com/webmutation))
- Resulting complete routing table has no 'tags "a,b"'  in last line [\#841](https://github.com/fabiolb/fabio/pull/841) ([hb9cwp](https://github.com/hb9cwp))
- fixes \#807 - changes map for grpc connections to be a string key [\#816](https://github.com/fabiolb/fabio/pull/816) ([nathanejohnson](https://github.com/nathanejohnson))
- Add command line flag to toggle required consistency on consul reads [\#811](https://github.com/fabiolb/fabio/pull/811) ([jeremycw](https://github.com/jeremycw))
- Issue 605 grpc host matching [\#632](https://github.com/fabiolb/fabio/pull/632) ([tommyalatalo](https://github.com/tommyalatalo))

## [v1.6.0](https://github.com/fabiolb/fabio/tree/v1.6.0) (2022-04-11)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.5.15...v1.6.0)

**Implemented enhancements:**

- Add support for influxdb metrics [\#253](https://github.com/fabiolb/fabio/issues/253)
- Support for prometheus [\#211](https://github.com/fabiolb/fabio/issues/211)
- Support dogstatd with tags [\#165](https://github.com/fabiolb/fabio/issues/165)
- Riemann metrics support [\#126](https://github.com/fabiolb/fabio/issues/126)
- Simple HTTP path prefix replacement [\#767](https://github.com/fabiolb/fabio/pull/767) ([JamesJJ](https://github.com/JamesJJ))

**Closed issues:**

- Consul Route updates very slow with large numbers of routes [\#865](https://github.com/fabiolb/fabio/issues/865)
- Restricting TLS versions [\#859](https://github.com/fabiolb/fabio/issues/859)
- \[admin/ui\] General updates [\#856](https://github.com/fabiolb/fabio/issues/856)
- \[question\] - Can Fabio listen on 80/tcp with Nomad [\#844](https://github.com/fabiolb/fabio/issues/844)
- Supporting requests in the form of /my-app/page1 [\#842](https://github.com/fabiolb/fabio/issues/842)
- Fabio Using Container IPs to create routes [\#839](https://github.com/fabiolb/fabio/issues/839)
- All my dynamic routes suddenly vanished! [\#837](https://github.com/fabiolb/fabio/issues/837)
- Fabio redirecting to /routes on own service [\#832](https://github.com/fabiolb/fabio/issues/832)
- docu fabio configure TLS/SSL\(HTTPS\) understanding problem [\#827](https://github.com/fabiolb/fabio/issues/827)
- Crash: \[FATAL\] 1.5.13. strconv.ParseUint: parsing ":1883;PROTO=TCP-DYNAMIC": invalid syntax [\#826](https://github.com/fabiolb/fabio/issues/826)
- strip doesn't work as expected on redirect [\#824](https://github.com/fabiolb/fabio/issues/824)
- Using Fabio with Consul over mTLS [\#820](https://github.com/fabiolb/fabio/issues/820)
- Switch to github actions [\#817](https://github.com/fabiolb/fabio/issues/817)
- Panic - httputil: ReverseProxy read error during body copy [\#814](https://github.com/fabiolb/fabio/issues/814)
- Support for Consul and Vault Namespaces [\#810](https://github.com/fabiolb/fabio/issues/810)
- grpc be closed when uninstall service target [\#807](https://github.com/fabiolb/fabio/issues/807)
- fabio binary filename for download [\#805](https://github.com/fabiolb/fabio/issues/805)
- Add arm64 arch [\#804](https://github.com/fabiolb/fabio/issues/804)
- Can Fabio to prefer one ethernet interface for proxy\_addr? [\#802](https://github.com/fabiolb/fabio/issues/802)
- TCP Dynamic Proxy route without specifying exact IP? [\#797](https://github.com/fabiolb/fabio/issues/797)
- \[Question\] What are opinions on allowing stale reads of Consul Catalog [\#764](https://github.com/fabiolb/fabio/issues/764)
- Simple HTTP path prefix replacement [\#691](https://github.com/fabiolb/fabio/issues/691)
- Does Fabio support multiple CS Stores per listener? [\#666](https://github.com/fabiolb/fabio/issues/666)
- \[Question\] Stats - Status code per service [\#371](https://github.com/fabiolb/fabio/issues/371)
- Statsd output is not good [\#327](https://github.com/fabiolb/fabio/issues/327)
- Send metrics to cloudwatch [\#326](https://github.com/fabiolb/fabio/issues/326)
- Mixing of HTTPS proxying and SNI+TCP on a single port [\#213](https://github.com/fabiolb/fabio/issues/213)

**Merged pull requests:**

- gofmt [\#870](https://github.com/fabiolb/fabio/pull/870) ([nathanejohnson](https://github.com/nathanejohnson))
- updating x/sys [\#869](https://github.com/fabiolb/fabio/pull/869) ([nathanejohnson](https://github.com/nathanejohnson))
- update go and alpine versions [\#868](https://github.com/fabiolb/fabio/pull/868) ([Netlims](https://github.com/Netlims))
- \#865 Move the route table sort into NewTable so that it only happens once. [\#867](https://github.com/fabiolb/fabio/pull/867) ([ddreier](https://github.com/ddreier))
- removing exclusion of arm64 mac build.  Fixes \#804 [\#866](https://github.com/fabiolb/fabio/pull/866) ([nathanejohnson](https://github.com/nathanejohnson))
- Fix example commands in registry.consul.kvpath [\#864](https://github.com/fabiolb/fabio/pull/864) ([blake](https://github.com/blake))
- Add IdleConnTimeout configurable for http transport [\#863](https://github.com/fabiolb/fabio/pull/863) ([aal89](https://github.com/aal89))
- admin/ui updates: [\#857](https://github.com/fabiolb/fabio/pull/857) ([dcarbone](https://github.com/dcarbone))
- Update 2 broken links in documentation [\#822](https://github.com/fabiolb/fabio/pull/822) ([mig4ng](https://github.com/mig4ng))
- Fix small typo in README.md [\#821](https://github.com/fabiolb/fabio/pull/821) ([mig4ng](https://github.com/mig4ng))
- Add support for github actions [\#819](https://github.com/fabiolb/fabio/pull/819) ([nathanejohnson](https://github.com/nathanejohnson))
- Remove golang toolchain name from release binary names [\#818](https://github.com/fabiolb/fabio/pull/818) ([nathanejohnson](https://github.com/nathanejohnson))
- we don't use Fabio [\#813](https://github.com/fabiolb/fabio/pull/813) ([hsmade](https://github.com/hsmade))
- Updating tcp dynamic proxy to match on routes that are port only [\#806](https://github.com/fabiolb/fabio/pull/806) ([nathanejohnson](https://github.com/nathanejohnson))
- Refactor metrics [\#476](https://github.com/fabiolb/fabio/pull/476) ([magiconair](https://github.com/magiconair))

## [v1.5.15](https://github.com/fabiolb/fabio/tree/v1.5.15) (2020-12-01)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.5.14...v1.5.15)

**Closed issues:**

- TCP Dynamic Proxy is not releasing ports from deregistered services [\#796](https://github.com/fabiolb/fabio/issues/796)
- How to configure log file output path [\#781](https://github.com/fabiolb/fabio/issues/781)

**Merged pull requests:**

- Updating the default GOGC to 100.  800 proves to be a bit insane. [\#803](https://github.com/fabiolb/fabio/pull/803) ([nathanejohnson](https://github.com/nathanejohnson))
- Stop dynamic TCP listener when upstream is no longer available [\#798](https://github.com/fabiolb/fabio/pull/798) ([fwkz](https://github.com/fwkz))
- Updating dependencies [\#794](https://github.com/fabiolb/fabio/pull/794) ([nathanejohnson](https://github.com/nathanejohnson))
- Update CHANGELOG.md [\#790](https://github.com/fabiolb/fabio/pull/790) ([stevenscg](https://github.com/stevenscg))

## [v1.5.14](https://github.com/fabiolb/fabio/tree/v1.5.14) (2020-09-09)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.5.13...v1.5.14)

**Fixed bugs:**

- %20 in route is causing route mismatch,  regression in 1.5.2 , works with 1.3.7 [\#347](https://github.com/fabiolb/fabio/issues/347)

**Closed issues:**

- matchingHostNoGlob sometimes returns incorrect matched host [\#786](https://github.com/fabiolb/fabio/issues/786)
- Add support for HTTPS+TCP+SNI on the same listener [\#783](https://github.com/fabiolb/fabio/issues/783)
- SIGTERM + Gracefully closing connections [\#782](https://github.com/fabiolb/fabio/issues/782)
- passing multiple routes via command line [\#776](https://github.com/fabiolb/fabio/issues/776)
- Master branch build failing with SECURITY ERROR [\#769](https://github.com/fabiolb/fabio/issues/769)
- How to disable client authentication for https? [\#765](https://github.com/fabiolb/fabio/issues/765)
- Must Access Control require RemoteAddr matching? [\#754](https://github.com/fabiolb/fabio/issues/754)
- Fabio Proxy \(localhost:9999\) Showing Blank White Screen [\#752](https://github.com/fabiolb/fabio/issues/752)
- Fabio 1.5.13 - no more "\[INFO\] Config updates" message in the logs [\#751](https://github.com/fabiolb/fabio/issues/751)
- Authentication issue. [\#743](https://github.com/fabiolb/fabio/issues/743)
- Connecting to HTTPS Upstream service. [\#738](https://github.com/fabiolb/fabio/issues/738)
- log.routes.format is broken with 1.5.13 [\#737](https://github.com/fabiolb/fabio/issues/737)
- Looking for a new maintainer [\#735](https://github.com/fabiolb/fabio/issues/735)
- GRPC Proxy + HTTP Proxy, both useable at the same time? [\#734](https://github.com/fabiolb/fabio/issues/734)
- Trace spans all have the same operation name [\#732](https://github.com/fabiolb/fabio/issues/732)
-  consul: Error fetching config from /fabio/config. Get  [\#729](https://github.com/fabiolb/fabio/issues/729)
- Very frequent 502 errors  [\#721](https://github.com/fabiolb/fabio/issues/721)
- Fabio decodes URL path parameters [\#486](https://github.com/fabiolb/fabio/issues/486)
- http proxy error context canceled [\#264](https://github.com/fabiolb/fabio/issues/264)

**Merged pull requests:**

- Fixing issue \#786 - matchingHostNoGlob sometimes returns incorrect host [\#787](https://github.com/fabiolb/fabio/pull/787) ([nathanejohnson](https://github.com/nathanejohnson))
- updating documentation for pending 1.5.14 release [\#785](https://github.com/fabiolb/fabio/pull/785) ([nathanejohnson](https://github.com/nathanejohnson))
- https+tcp+sni listener support [\#784](https://github.com/fabiolb/fabio/pull/784) ([nathanejohnson](https://github.com/nathanejohnson))
- chore: fix typo in comments [\#775](https://github.com/fabiolb/fabio/pull/775) ([josgraha](https://github.com/josgraha))
- \(docs\): fixed small error [\#774](https://github.com/fabiolb/fabio/pull/774) ([0xflotus](https://github.com/0xflotus))
- Preserve table state by storing buffer table in fixed strings. [\#749](https://github.com/fabiolb/fabio/pull/749) ([aaronhurt](https://github.com/aaronhurt))
- only deploy once per build [\#747](https://github.com/fabiolb/fabio/pull/747) ([aaronhurt](https://github.com/aaronhurt))
- switch to github pages for doc hosting [\#746](https://github.com/fabiolb/fabio/pull/746) ([aaronhurt](https://github.com/aaronhurt))
- minor transition updates and small fixes [\#745](https://github.com/fabiolb/fabio/pull/745) ([aaronhurt](https://github.com/aaronhurt))
- switch back to travis CI [\#744](https://github.com/fabiolb/fabio/pull/744) ([nathanejohnson](https://github.com/nathanejohnson))
- follow hugo best practices [\#742](https://github.com/fabiolb/fabio/pull/742) ([aaronhurt](https://github.com/aaronhurt))
- Documentation updates for project transition. [\#740](https://github.com/fabiolb/fabio/pull/740) ([aaronhurt](https://github.com/aaronhurt))
- Fix infinite buffering of SSE responses when gzip is enabled [\#739](https://github.com/fabiolb/fabio/pull/739) ([ctlajoie](https://github.com/ctlajoie))
- Add missing \<svc\> entry to example route [\#733](https://github.com/fabiolb/fabio/pull/733) ([BenjaminHerbert](https://github.com/BenjaminHerbert))
- minor fixups [\#731](https://github.com/fabiolb/fabio/pull/731) ([aaronhurt](https://github.com/aaronhurt))
- fix tests [\#730](https://github.com/fabiolb/fabio/pull/730) ([aaronhurt](https://github.com/aaronhurt))
- Add HTTP method and path to trace span operation name [\#715](https://github.com/fabiolb/fabio/pull/715) ([hobochili](https://github.com/hobochili))
- Deprecate deregisterCriticalServiceAfter option [\#674](https://github.com/fabiolb/fabio/pull/674) ([pschultz](https://github.com/pschultz))
- Issue 647 NormalizeHost [\#648](https://github.com/fabiolb/fabio/pull/648) ([murphymj25](https://github.com/murphymj25))
- Handle context canceled errors + better http proxy error handling [\#644](https://github.com/fabiolb/fabio/pull/644) ([danlsgiga](https://github.com/danlsgiga))
- Added idleTimout to config and to serve.go HTTP server [\#635](https://github.com/fabiolb/fabio/pull/635) ([galen0624](https://github.com/galen0624))
- Issue 613 tcp dynamic [\#626](https://github.com/fabiolb/fabio/pull/626) ([murphymj25](https://github.com/murphymj25))
- Issue 554 Added compiled glob matching using LRU Cache [\#615](https://github.com/fabiolb/fabio/pull/615) ([galen0624](https://github.com/galen0624))
- Issue 558 - Add Polling Interval From Fabio to Consul to Fabio Config [\#572](https://github.com/fabiolb/fabio/pull/572) ([galen0624](https://github.com/galen0624))
- Feat: Pass encoded characters in path unchanged [\#489](https://github.com/fabiolb/fabio/pull/489) ([valentin-krasontovitsch](https://github.com/valentin-krasontovitsch))

## [v1.5.13](https://github.com/fabiolb/fabio/tree/v1.5.13) (2019-11-18)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.5.12...v1.5.13)

**Closed issues:**

- Fabio 1.5.12 - panic: runtime error: invalid memory address or nil pointer dereference [\#719](https://github.com/fabiolb/fabio/issues/719)
- Question: Load balancing WebSocket connections [\#718](https://github.com/fabiolb/fabio/issues/718)
- Question: resources \(css, js files\) by multiple sites [\#717](https://github.com/fabiolb/fabio/issues/717)
- Fabio UI not displaying when hit on a DNS name [\#712](https://github.com/fabiolb/fabio/issues/712)
- Unable to route to websites [\#676](https://github.com/fabiolb/fabio/issues/676)
- Websocket proxy timeouts [\#518](https://github.com/fabiolb/fabio/issues/518)

**Merged pull requests:**

- fix nil-pointer dereference in detailed config log [\#720](https://github.com/fabiolb/fabio/pull/720) ([pschultz](https://github.com/pschultz))
- Safely handle missing cert from Vault KV store [\#710](https://github.com/fabiolb/fabio/pull/710) ([dradtke](https://github.com/dradtke))

## [v1.5.12](https://github.com/fabiolb/fabio/tree/v1.5.12) (2019-10-11)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.5.11...v1.5.12)

**Implemented enhancements:**

- docker swarm some times register eth0 other eth1 [\#652](https://github.com/fabiolb/fabio/issues/652)
- config: let registry.consul.register.addr default to ui.addr [\#658](https://github.com/fabiolb/fabio/pull/658) ([pschultz](https://github.com/pschultz))
- fix exit status code [\#637](https://github.com/fabiolb/fabio/pull/637) ([ianic](https://github.com/ianic))

**Closed issues:**

- Example of Vault KV clientca option? [\#703](https://github.com/fabiolb/fabio/issues/703)
- tcp proxy not work [\#702](https://github.com/fabiolb/fabio/issues/702)
- urlprefix-:3306 proto=tcp   not work [\#701](https://github.com/fabiolb/fabio/issues/701)
- https proxy not work [\#700](https://github.com/fabiolb/fabio/issues/700)
- the http port is 9999 ,the https port is what? [\#699](https://github.com/fabiolb/fabio/issues/699)
- TCP proxy log filled with i/o timeout [\#696](https://github.com/fabiolb/fabio/issues/696)
- urlprefix-zzz.xxx.com/api  not work [\#693](https://github.com/fabiolb/fabio/issues/693)
- Fabio/Consul route integration [\#689](https://github.com/fabiolb/fabio/issues/689)
- Unable to route. [\#680](https://github.com/fabiolb/fabio/issues/680)
- Fabio 100% CPU usage due to logging [\#673](https://github.com/fabiolb/fabio/issues/673)
- Authorization header leaking to the backend. [\#671](https://github.com/fabiolb/fabio/issues/671)
- X-Request-Start header [\#670](https://github.com/fabiolb/fabio/issues/670)
- fabio service entries may stay in Consul on dirty exit [\#663](https://github.com/fabiolb/fabio/issues/663)
- Can fabio route request by request body [\#661](https://github.com/fabiolb/fabio/issues/661)
- Wrong reported HealthCheck-URI using custom -proxy.addr & -ui.addr [\#657](https://github.com/fabiolb/fabio/issues/657)
- Clarify documentation HTTP Redirects [\#656](https://github.com/fabiolb/fabio/issues/656)
- tcp access control doesn't work [\#651](https://github.com/fabiolb/fabio/issues/651)
- Crash on start of watchBackend\(\) [\#650](https://github.com/fabiolb/fabio/issues/650)
- Remove third-party cookie and script requirements from frontend [\#642](https://github.com/fabiolb/fabio/issues/642)
- Build should use included vendor directory with modules [\#638](https://github.com/fabiolb/fabio/issues/638)
- Route table UI is broken [\#628](https://github.com/fabiolb/fabio/issues/628)
- Possible Memory Leak in WatchBackend [\#595](https://github.com/fabiolb/fabio/issues/595)
- Release date for 1.5.11 [\#592](https://github.com/fabiolb/fabio/issues/592)
- Fabio and Vault Token Issues [\#523](https://github.com/fabiolb/fabio/issues/523)
- UI broken where no internet access. [\#502](https://github.com/fabiolb/fabio/issues/502)
- make log compatible with the syslog protocol [\#397](https://github.com/fabiolb/fabio/issues/397)

**Merged pull requests:**

- Add Vault example to the traffic shaping section. [\#677](https://github.com/fabiolb/fabio/pull/677) ([jrasell](https://github.com/jrasell))
- Fix matching priority for host:port tuples [\#675](https://github.com/fabiolb/fabio/pull/675) ([pschultz](https://github.com/pschultz))
- Add config option to use 128 bit trace IDs [\#669](https://github.com/fabiolb/fabio/pull/669) ([gfloyd](https://github.com/gfloyd))
- register: clean-up fabio service entries in Consul on dirty exit [\#664](https://github.com/fabiolb/fabio/pull/664) ([pires](https://github.com/pires))
- Fix SSE by implementing Flusher in responseWriter wrapper [\#655](https://github.com/fabiolb/fabio/pull/655) ([gfloyd](https://github.com/gfloyd))
- Use go-sockaddr to parse address strings [\#653](https://github.com/fabiolb/fabio/pull/653) ([aaronhurt](https://github.com/aaronhurt))
- ensure absolute path after strip to maintain rfc complaince [\#645](https://github.com/fabiolb/fabio/pull/645) ([aaronhurt](https://github.com/aaronhurt))
- Bundle UI assets [\#643](https://github.com/fabiolb/fabio/pull/643) ([pschultz](https://github.com/pschultz))
- ui: Remove duplicate destination column [\#641](https://github.com/fabiolb/fabio/pull/641) ([pschultz](https://github.com/pschultz))
- use vendor directory when building - fixes \#638 [\#639](https://github.com/fabiolb/fabio/pull/639) ([aaronhurt](https://github.com/aaronhurt))
- Issue 595 watchbackend [\#629](https://github.com/fabiolb/fabio/pull/629) ([murphymj25](https://github.com/murphymj25))
- added support for profile/tracing [\#624](https://github.com/fabiolb/fabio/pull/624) ([galen0624](https://github.com/galen0624))

## [v1.5.11](https://github.com/fabiolb/fabio/tree/v1.5.11) (2019-04-09)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.5.11-wrong...v1.5.11)

**Implemented enhancements:**

- Proxy protocol support fo outgoing connections [\#191](https://github.com/fabiolb/fabio/issues/191)

**Closed issues:**

- Consul blocking queries should be rate limited to avoid spiking loads on server [\#627](https://github.com/fabiolb/fabio/issues/627)
- This seems to be a recursive func call.  Is this correct? [\#625](https://github.com/fabiolb/fabio/issues/625)
- Bug in consul 1.4.3 [\#616](https://github.com/fabiolb/fabio/issues/616)
- \[question\] Release date for 1.5.11 [\#601](https://github.com/fabiolb/fabio/issues/601)
- Sidebar of the website is a little off [\#599](https://github.com/fabiolb/fabio/issues/599)
- wrong use  function strings.HasPrefix\(\) in file passsing.go [\#545](https://github.com/fabiolb/fabio/issues/545)
- best way to bypass fabio consul integration? [\#437](https://github.com/fabiolb/fabio/issues/437)

**Merged pull requests:**

- Issue 611 Added Custom API Driven Back end [\#614](https://github.com/fabiolb/fabio/pull/614) ([galen0624](https://github.com/galen0624))
- Improved basic auth htpasswd file refresh \#604 [\#610](https://github.com/fabiolb/fabio/pull/610) ([mfuterko](https://github.com/mfuterko))
- Address \#545 - wrong use function strings.HasPrefix [\#607](https://github.com/fabiolb/fabio/pull/607) ([mfuterko](https://github.com/mfuterko))
- docs: fix layout without JS enabled [\#606](https://github.com/fabiolb/fabio/pull/606) ([pschultz](https://github.com/pschultz))
- Implement basic auth htpasswd file refresh [\#604](https://github.com/fabiolb/fabio/pull/604) ([mfuterko](https://github.com/mfuterko))
- added support for Consul TLS transport [\#602](https://github.com/fabiolb/fabio/pull/602) ([sev3ryn](https://github.com/sev3ryn))
- Proxy protocol on outbound tcp, tcp+sni and tcp with tls connection [\#598](https://github.com/fabiolb/fabio/pull/598) ([mfuterko](https://github.com/mfuterko))

## [v1.5.11-wrong](https://github.com/fabiolb/fabio/tree/v1.5.11-wrong) (2019-02-25)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.5.10...v1.5.11-wrong)

**Implemented enhancements:**

- Basic authentication on routes [\#166](https://github.com/fabiolb/fabio/issues/166)

**Fixed bugs:**

- TCP proxy broken since v1.5.8 [\#524](https://github.com/fabiolb/fabio/issues/524)

**Closed issues:**

- Fabio's routing table empty. Consul indicates registered services with urlprefix- tags [\#589](https://github.com/fabiolb/fabio/issues/589)
- HTTP 502 response half of the time [\#584](https://github.com/fabiolb/fabio/issues/584)
- tcp+sni route with allow=ip:something does not seem to work [\#576](https://github.com/fabiolb/fabio/issues/576)
- Passing args to fabio in nomad task. [\#567](https://github.com/fabiolb/fabio/issues/567)
- Change Log entry update  [\#562](https://github.com/fabiolb/fabio/issues/562)
- Release date for 1.5.10? [\#560](https://github.com/fabiolb/fabio/issues/560)
- Route updates are delayed with large number of services  [\#558](https://github.com/fabiolb/fabio/issues/558)
- could the source and destination be clickable in the ui? [\#508](https://github.com/fabiolb/fabio/issues/508)
- Support for opentracing [\#429](https://github.com/fabiolb/fabio/issues/429)
- Case-insensitive path matching [\#35](https://github.com/fabiolb/fabio/issues/35)

**Merged pull requests:**

- ui: Fix XSS vulnerability [\#588](https://github.com/fabiolb/fabio/pull/588) ([pschultz](https://github.com/pschultz))
- make Dest column into clickable links [\#587](https://github.com/fabiolb/fabio/pull/587) ([kneufeld](https://github.com/kneufeld))
- update documentation around the changes to PROXY protocol [\#583](https://github.com/fabiolb/fabio/pull/583) ([aaronhurt](https://github.com/aaronhurt))
- address concerns raised while troubleshooting \#524 [\#581](https://github.com/fabiolb/fabio/pull/581) ([aaronhurt](https://github.com/aaronhurt))
- fix ip access rules within tcp proxy - fixes \#576 [\#577](https://github.com/fabiolb/fabio/pull/577) ([aaronhurt](https://github.com/aaronhurt))
- Add GRPC proxy support [\#575](https://github.com/fabiolb/fabio/pull/575) ([andyroyle](https://github.com/andyroyle))
- metrics.circonus: Add support for circonus.submissionurl [\#574](https://github.com/fabiolb/fabio/pull/574) ([stack72](https://github.com/stack72))
- add http-basic auth reading from a file [\#573](https://github.com/fabiolb/fabio/pull/573) ([andyroyle](https://github.com/andyroyle))
- update consul to v1.4.0 - fixes \#569 [\#571](https://github.com/fabiolb/fabio/pull/571) ([aaronhurt](https://github.com/aaronhurt))
- add faq to address \#490 [\#568](https://github.com/fabiolb/fabio/pull/568) ([aaronhurt](https://github.com/aaronhurt))
- Update go.mod for \#472 [\#565](https://github.com/fabiolb/fabio/pull/565) ([magiconair](https://github.com/magiconair))
- Refactor consul service monitor [\#564](https://github.com/fabiolb/fabio/pull/564) ([magiconair](https://github.com/magiconair))
- issue 562 update change log glob.matching.disabled [\#563](https://github.com/fabiolb/fabio/pull/563) ([galen0624](https://github.com/galen0624))
- Added new case insensitive matcher [\#553](https://github.com/fabiolb/fabio/pull/553) ([herbrandson](https://github.com/herbrandson))
- \[Docs\] Delete duplicate 'Path Stripping' page [\#537](https://github.com/fabiolb/fabio/pull/537) ([rkettelerij](https://github.com/rkettelerij))
- \#429 issue - OpenTrace zipKin Support  [\#472](https://github.com/fabiolb/fabio/pull/472) ([galen0624](https://github.com/galen0624))

## [v1.5.10](https://github.com/fabiolb/fabio/tree/v1.5.10) (2018-10-25)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.5.9...v1.5.10)

**Fixed bugs:**

- Wrong route for multiple matching host glob patterns [\#506](https://github.com/fabiolb/fabio/issues/506)

**Closed issues:**

- Fabio forcing response header keys upper case [\#552](https://github.com/fabiolb/fabio/issues/552)
- Multiple fabio instances load balancing different set of services. [\#551](https://github.com/fabiolb/fabio/issues/551)
- Without Consul how can i use Fabio? [\#549](https://github.com/fabiolb/fabio/issues/549)
- Performance issue - Glob matching with large number of services in consul [\#548](https://github.com/fabiolb/fabio/issues/548)
- Ignore host case when adding and matching routes [\#542](https://github.com/fabiolb/fabio/issues/542)
- allow redirect host to be empty [\#533](https://github.com/fabiolb/fabio/issues/533)
- Expose Fabio metrics via Prometheus [\#532](https://github.com/fabiolb/fabio/issues/532)
- Memory leak in go-metrics library [\#530](https://github.com/fabiolb/fabio/issues/530)
- ability to remove headers from the request [\#528](https://github.com/fabiolb/fabio/issues/528)
- urlprefix- does not work properly [\#527](https://github.com/fabiolb/fabio/issues/527)
- Problem geting fabio routing to its own ui [\#525](https://github.com/fabiolb/fabio/issues/525)
- Redirection to default back-end if route not exists [\#521](https://github.com/fabiolb/fabio/issues/521)
- Forwarding Uri tag in original request to endpoint [\#519](https://github.com/fabiolb/fabio/issues/519)
- Fabio - Manual overrides [\#515](https://github.com/fabiolb/fabio/issues/515)
- If consul is behind an ELB with a set timeout, and the connection is timed out by the ELB, subsequent requests from fabio fail [\#513](https://github.com/fabiolb/fabio/issues/513)
- Fabio instantly delete route, whereas health check is passing [\#512](https://github.com/fabiolb/fabio/issues/512)
- Would it be possible to configure Fabio to watch services with warning state? [\#509](https://github.com/fabiolb/fabio/issues/509)
- Headers passed through fabio are modified [\#505](https://github.com/fabiolb/fabio/issues/505)
- Fabio -\> HTTPS -\> Service ? [\#503](https://github.com/fabiolb/fabio/issues/503)
- Tls + sni support for non http traffic?  [\#499](https://github.com/fabiolb/fabio/issues/499)
- Static routes in fabio.properties [\#498](https://github.com/fabiolb/fabio/issues/498)
- Tests fail with consul \> 1.0.6 and vault \> 0.9.6 [\#494](https://github.com/fabiolb/fabio/issues/494)
- Question: wildcard hostname support [\#491](https://github.com/fabiolb/fabio/issues/491)
- Fabio confi help with multiple proto [\#490](https://github.com/fabiolb/fabio/issues/490)
- Add support for Vault 0.10 KV v2 [\#483](https://github.com/fabiolb/fabio/issues/483)
- Support "standard" Consul envvars [\#277](https://github.com/fabiolb/fabio/issues/277)
- Support Consul TLS [\#276](https://github.com/fabiolb/fabio/issues/276)

**Merged pull requests:**

- Issue \#548 added enable/disable glob matching [\#550](https://github.com/fabiolb/fabio/pull/550) ([galen0624](https://github.com/galen0624))
- Correct the access control feature documentation page [\#546](https://github.com/fabiolb/fabio/pull/546) ([msvbhat](https://github.com/msvbhat))
- Add $host pseudo variable [\#544](https://github.com/fabiolb/fabio/pull/544) ([holtwilkins](https://github.com/holtwilkins))
- compare host using lowercase [\#543](https://github.com/fabiolb/fabio/pull/543) ([shantanugadgil](https://github.com/shantanugadgil))
- Issue \#530 - Vendored in updated go-metrics package [\#535](https://github.com/fabiolb/fabio/pull/535) ([galen0624](https://github.com/galen0624))
- Add setting to flush fabio buffer regardless headers [\#531](https://github.com/fabiolb/fabio/pull/531) ([samm-git](https://github.com/samm-git))
- Update README.md [\#510](https://github.com/fabiolb/fabio/pull/510) ([kuskmen](https://github.com/kuskmen))
- Issue \#506: reverse domain names before sorting [\#507](https://github.com/fabiolb/fabio/pull/507) ([magiconair](https://github.com/magiconair))
- Fix changelog link in docs footer [\#500](https://github.com/fabiolb/fabio/pull/500) ([xmikus01](https://github.com/xmikus01))
- Make tests compatible with Vault 0.10 [\#497](https://github.com/fabiolb/fabio/pull/497) ([pschultz](https://github.com/pschultz))
- Delete an unused global variable logOutput [\#495](https://github.com/fabiolb/fabio/pull/495) ([gua-pian](https://github.com/gua-pian))
- Add fastcgi handler [\#435](https://github.com/fabiolb/fabio/pull/435) ([Gufran](https://github.com/Gufran))

## [v1.5.9](https://github.com/fabiolb/fabio/tree/v1.5.9) (2018-05-16)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.5.8...v1.5.9)

**Closed issues:**

- UI is broken from  versions =\> 1.7 [\#487](https://github.com/fabiolb/fabio/issues/487)
- Building master fails [\#482](https://github.com/fabiolb/fabio/issues/482)
- '-registry.consul.register.enabled' does not seem to be respected [\#467](https://github.com/fabiolb/fabio/issues/467)
- Access logging fails in combination with proxy gzipping [\#460](https://github.com/fabiolb/fabio/issues/460)
- glob matching improvements [\#452](https://github.com/fabiolb/fabio/issues/452)
- Add route based on x-forwarded-port header [\#450](https://github.com/fabiolb/fabio/issues/450)
- Redirect http to https on the same destination [\#448](https://github.com/fabiolb/fabio/issues/448)
- WebSocket Upgrade not sending Response [\#447](https://github.com/fabiolb/fabio/issues/447)
- Fabio does not remove service when one of the registered health-checks fail [\#427](https://github.com/fabiolb/fabio/issues/427)
- Fabio routing to wrong back end [\#421](https://github.com/fabiolb/fabio/issues/421)
- \[feature\]: proxy route option [\#356](https://github.com/fabiolb/fabio/issues/356)

**Merged pull requests:**

- Resetting read deadline [\#492](https://github.com/fabiolb/fabio/pull/492) ([craigday](https://github.com/craigday))
- Issue \#466: make redirect code more robust [\#477](https://github.com/fabiolb/fabio/pull/477) ([magiconair](https://github.com/magiconair))
- fix contributors link [\#475](https://github.com/fabiolb/fabio/pull/475) ([aaronhurt](https://github.com/aaronhurt))
- ws close on failed handshake \(\#421\) [\#474](https://github.com/fabiolb/fabio/pull/474) ([magiconair](https://github.com/magiconair))
- Issue \#460: Fix access logging when gzip is enabled [\#470](https://github.com/fabiolb/fabio/pull/470) ([magiconair](https://github.com/magiconair))
- Fix the regex of the example proxy.gzip.contenttype [\#468](https://github.com/fabiolb/fabio/pull/468) ([tino](https://github.com/tino))
- Check upstream X-Forwarded-Proto prior to redirect [\#466](https://github.com/fabiolb/fabio/pull/466) ([aaronhurt](https://github.com/aaronhurt))
- Fix certificate stores doc path [\#458](https://github.com/fabiolb/fabio/pull/458) ([eldondev](https://github.com/eldondev))
- Add new & improved glob matcher [\#457](https://github.com/fabiolb/fabio/pull/457) ([sharbov](https://github.com/sharbov))
- handle indeterminate length proxy chains - fixes \#449 [\#453](https://github.com/fabiolb/fabio/pull/453) ([aaronhurt](https://github.com/aaronhurt))
- Update link for Websockets [\#446](https://github.com/fabiolb/fabio/pull/446) ([a2ar](https://github.com/a2ar))
- "strict" health-checking \(\#427\) [\#428](https://github.com/fabiolb/fabio/pull/428) ([systemfreund](https://github.com/systemfreund))

## [v1.5.8](https://github.com/fabiolb/fabio/tree/v1.5.8) (2018-02-18)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.5.7...v1.5.8)

**Closed issues:**

- TCP Proxying SSH connections [\#445](https://github.com/fabiolb/fabio/issues/445)
- route add ... opts "proto=tcp+sni" ?? [\#444](https://github.com/fabiolb/fabio/issues/444)
- Wildcard registeration issues [\#440](https://github.com/fabiolb/fabio/issues/440)
- Feature Request: IP Whitelisting [\#439](https://github.com/fabiolb/fabio/issues/439)
- NoRouteHTMLPath not rendering HTML page [\#438](https://github.com/fabiolb/fabio/issues/438)

**Merged pull requests:**

- ignore fabio.exe [\#443](https://github.com/fabiolb/fabio/pull/443) ([aaronhurt](https://github.com/aaronhurt))
- Issue \#438: Do not add separators for NoRouteHTML page [\#441](https://github.com/fabiolb/fabio/pull/441) ([magiconair](https://github.com/magiconair))
- Add option to allow Fabio to register frontend services in Consul on behalf of user services [\#426](https://github.com/fabiolb/fabio/pull/426) ([rileyje](https://github.com/rileyje))
- TCP+SNI support arbitrary large Client Hello [\#423](https://github.com/fabiolb/fabio/pull/423) ([DanSipola](https://github.com/DanSipola))

## [v1.5.7](https://github.com/fabiolb/fabio/tree/v1.5.7) (2018-02-06)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.5.6...v1.5.7)

**Closed issues:**

- VaultPKI tests fail with go1.10rc1 [\#434](https://github.com/fabiolb/fabio/issues/434)
- Ensure that proxy.noroutestatus has three digits [\#433](https://github.com/fabiolb/fabio/issues/433)
- Vault PKI documentation and Fabio version [\#430](https://github.com/fabiolb/fabio/issues/430)
- configure equivalent  of nginx client\_max\_body\_size [\#422](https://github.com/fabiolb/fabio/issues/422)
- \[question\] Newbie question: where to place urlpref-host/path? [\#419](https://github.com/fabiolb/fabio/issues/419)
- Static / Manual routes management via API [\#396](https://github.com/fabiolb/fabio/issues/396)
- Warn if fabio is run as root [\#369](https://github.com/fabiolb/fabio/issues/369)

**Merged pull requests:**

- Activating Open Collective [\#432](https://github.com/fabiolb/fabio/pull/432) ([monkeywithacupcake](https://github.com/monkeywithacupcake))
- fix small typo [\#431](https://github.com/fabiolb/fabio/pull/431) ([aaronhurt](https://github.com/aaronhurt))
- Add support for HSTS response headers and provide method for adding additional response headers [\#425](https://github.com/fabiolb/fabio/pull/425) ([aaronhurt](https://github.com/aaronhurt))
- Fix maxconn documentation [\#420](https://github.com/fabiolb/fabio/pull/420) ([slobo](https://github.com/slobo))
- treat registry.consul.kvpath as prefix [\#417](https://github.com/fabiolb/fabio/pull/417) ([magiconair](https://github.com/magiconair))
- Issue \#369: Do not allow to run fabio as root [\#377](https://github.com/fabiolb/fabio/pull/377) ([magiconair](https://github.com/magiconair))

## [v1.5.6](https://github.com/fabiolb/fabio/tree/v1.5.6) (2018-01-05)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.5.5...v1.5.6)

**Closed issues:**

- Excessive consul logging [\#408](https://github.com/fabiolb/fabio/issues/408)
- Build new website [\#405](https://github.com/fabiolb/fabio/issues/405)
- \[bug?\] Fabio uses "global" Consul ServiceID's [\#383](https://github.com/fabiolb/fabio/issues/383)

**Merged pull requests:**

- Issue \#408: log consul state changes as DEBUG [\#418](https://github.com/fabiolb/fabio/pull/418) ([magiconair](https://github.com/magiconair))
- Actually respect -version option [\#415](https://github.com/fabiolb/fabio/pull/415) ([pschultz](https://github.com/pschultz))
- Identify services using both the ID and the Node [\#414](https://github.com/fabiolb/fabio/pull/414) ([alvaroaleman](https://github.com/alvaroaleman))

## [v1.5.5](https://github.com/fabiolb/fabio/tree/v1.5.5) (2017-12-20)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.5.4...v1.5.5)

**Implemented enhancements:**

- Support custom 404/503 error pages [\#56](https://github.com/fabiolb/fabio/issues/56)

**Closed issues:**

- Fabio for task/container/service load balancing on amazon ecs with consul and registrator.  [\#402](https://github.com/fabiolb/fabio/issues/402)

**Merged pull requests:**

- Implement custom noroute html response [\#398](https://github.com/fabiolb/fabio/pull/398) ([tino](https://github.com/tino))

## [v1.5.4](https://github.com/fabiolb/fabio/tree/v1.5.4) (2017-12-10)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.5.3...v1.5.4)

**Implemented enhancements:**

- Differentiate "URL Unavailable/503" and "URL Not Found/404" [\#214](https://github.com/fabiolb/fabio/issues/214)

**Fixed bugs:**

- opts with host= with multiple routes does not work as expected [\#385](https://github.com/fabiolb/fabio/issues/385)

**Closed issues:**

- Fabio is not handling SIGHUP \(HUP\) signal properly - it dies [\#400](https://github.com/fabiolb/fabio/issues/400)
- Typo in manual overrides stops Fabio from updating routes [\#399](https://github.com/fabiolb/fabio/issues/399)
- route precendence  [\#389](https://github.com/fabiolb/fabio/issues/389)
- how to connect consul cluster [\#386](https://github.com/fabiolb/fabio/issues/386)
- Allow comments in manual overrides [\#379](https://github.com/fabiolb/fabio/issues/379)
- Domain or protocol redirection [\#87](https://github.com/fabiolb/fabio/issues/87)
- Should rewrite the Host Header  [\#75](https://github.com/fabiolb/fabio/issues/75)

**Merged pull requests:**

- Issue \#400: ignore SIGHUP [\#403](https://github.com/fabiolb/fabio/pull/403) ([magiconair](https://github.com/magiconair))
- Issue \#389: match exact host before glob matches [\#390](https://github.com/fabiolb/fabio/pull/390) ([magiconair](https://github.com/magiconair))
- Issue \#385: attach options to target instead of route [\#388](https://github.com/fabiolb/fabio/pull/388) ([magiconair](https://github.com/magiconair))
- Fix various minor things [\#382](https://github.com/fabiolb/fabio/pull/382) ([antham](https://github.com/antham))
- Remove unused variable [\#381](https://github.com/fabiolb/fabio/pull/381) ([antham](https://github.com/antham))
- Now setting the X-Forwarded-Host header if not present. Add matching … [\#380](https://github.com/fabiolb/fabio/pull/380) ([LeReverandNox](https://github.com/LeReverandNox))

## [v1.5.3](https://github.com/fabiolb/fabio/tree/v1.5.3) (2017-11-03)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.5.2...v1.5.3)

**Implemented enhancements:**

- Drop privileges after start [\#195](https://github.com/fabiolb/fabio/issues/195)
- support for adding CORS headers? [\#110](https://github.com/fabiolb/fabio/issues/110)

**Closed issues:**

- host=www.mydomain.com not working [\#375](https://github.com/fabiolb/fabio/issues/375)
- Wildcards in routing path [\#374](https://github.com/fabiolb/fabio/issues/374)
- Questions/issues in using overrides [\#372](https://github.com/fabiolb/fabio/issues/372)
- nodes and services in maintenance can cause excessive logging [\#367](https://github.com/fabiolb/fabio/issues/367)
- Support fabio.properties in Consul KV store [\#365](https://github.com/fabiolb/fabio/issues/365)
- Fabio fails to strip the prefix if the url prefix does not start with the strip option value [\#363](https://github.com/fabiolb/fabio/issues/363)
- More than one fabio instance decreases system performance. [\#361](https://github.com/fabiolb/fabio/issues/361)
- Documentation of the available metrics? [\#360](https://github.com/fabiolb/fabio/issues/360)
- select color scheme from config to distinguish environments [\#359](https://github.com/fabiolb/fabio/issues/359)
- \[Feature request\]: TCP Proxy support different incoming and outbound ports [\#353](https://github.com/fabiolb/fabio/issues/353)
- hgfiii [\#351](https://github.com/fabiolb/fabio/issues/351)
- statsd - unable to parse line - gf metric [\#350](https://github.com/fabiolb/fabio/issues/350)
- Possibility for Docker Image to pass Consul IP and Port as Variable? [\#346](https://github.com/fabiolb/fabio/issues/346)
- Ways to have log verbosity [\#345](https://github.com/fabiolb/fabio/issues/345)
- Cant disable consul register with -registry.consul.register.enabled=false [\#342](https://github.com/fabiolb/fabio/issues/342)
- Glob Matcher is not working for me [\#341](https://github.com/fabiolb/fabio/issues/341)
- Strip option has no effect for websockets [\#330](https://github.com/fabiolb/fabio/issues/330)
- access logging is not right [\#322](https://github.com/fabiolb/fabio/issues/322)
- FATAL error when metrics cannot be delivered [\#320](https://github.com/fabiolb/fabio/issues/320)
- http: proxy error: context canceled [\#318](https://github.com/fabiolb/fabio/issues/318)
- /api/routes intermittently returns null. [\#316](https://github.com/fabiolb/fabio/issues/316)
- what is the tcp writeTimeout? [\#307](https://github.com/fabiolb/fabio/issues/307)

**Merged pull requests:**

- Issue \#375: set host header when host option is set [\#376](https://github.com/fabiolb/fabio/pull/376) ([magiconair](https://github.com/magiconair))

## [v1.5.2](https://github.com/fabiolb/fabio/tree/v1.5.2) (2017-07-24)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.5.1...v1.5.2)

**Implemented enhancements:**

- Auto-generated Vault certs [\#135](https://github.com/fabiolb/fabio/issues/135)

**Closed issues:**

- not able to acces the service via fabio. [\#319](https://github.com/fabiolb/fabio/issues/319)

**Merged pull requests:**

- Fix memory leak in tcp proxy [\#321](https://github.com/fabiolb/fabio/pull/321) ([Crypto89](https://github.com/Crypto89))

## [v1.5.1](https://github.com/fabiolb/fabio/tree/v1.5.1) (2017-07-06)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.5.0...v1.5.1)

**Implemented enhancements:**

- Feature: Allow weight tag in Consul [\#42](https://github.com/fabiolb/fabio/issues/42)

**Fixed bugs:**

- 1.5.0 config compatibility problem  [\#305](https://github.com/fabiolb/fabio/issues/305)

**Closed issues:**

- Multiple urlprefix [\#317](https://github.com/fabiolb/fabio/issues/317)
- Add metrics for TCP and TCP+SNI proxy [\#306](https://github.com/fabiolb/fabio/issues/306)
- How to configure TCP correctly \(proxy.addr, ...\) [\#283](https://github.com/fabiolb/fabio/issues/283)
- Add parameter to vault token renewal [\#274](https://github.com/fabiolb/fabio/issues/274)

**Merged pull requests:**

- Issue \#274: Avoid premature vault token renewals [\#314](https://github.com/fabiolb/fabio/pull/314) ([pschultz](https://github.com/pschultz))
- Make tests work with vault 0.7.x [\#313](https://github.com/fabiolb/fabio/pull/313) ([pschultz](https://github.com/pschultz))
- Fix syntax highlighting in README [\#311](https://github.com/fabiolb/fabio/pull/311) ([agis](https://github.com/agis))

## [v1.5.0](https://github.com/fabiolb/fabio/tree/v1.5.0) (2017-06-07)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.4.4...v1.5.0)

**Implemented enhancements:**

- X-Forwarded-Prefix header support [\#304](https://github.com/fabiolb/fabio/issues/304)
- read only web ui [\#302](https://github.com/fabiolb/fabio/issues/302)
- Sync X-Forwarded-Proto and Forwarded header when possible [\#296](https://github.com/fabiolb/fabio/issues/296)
- Using upstream hostname for request [\#294](https://github.com/fabiolb/fabio/issues/294)
- Add profiling support [\#290](https://github.com/fabiolb/fabio/issues/290)
- TLS and Connection information through headers [\#280](https://github.com/fabiolb/fabio/issues/280)
- Support TLS/Ciphersuite configuration options [\#249](https://github.com/fabiolb/fabio/issues/249)

**Fixed bugs:**

- Support gzip compression for websockets [\#300](https://github.com/fabiolb/fabio/issues/300)

**Closed issues:**

- Example of proxy.gzip.contenttype configuration [\#299](https://github.com/fabiolb/fabio/issues/299)
- Compatibility with 1.8 [\#297](https://github.com/fabiolb/fabio/issues/297)
- cert file names and path= not working as documented [\#293](https://github.com/fabiolb/fabio/issues/293)
- Multiple SSL certs for same listener [\#291](https://github.com/fabiolb/fabio/issues/291)
- HTTPProxy cannot be aware of timeout of waiting response [\#288](https://github.com/fabiolb/fabio/issues/288)
- websockets failing with 500 response - running rancher behind fabio [\#133](https://github.com/fabiolb/fabio/issues/133)

**Merged pull requests:**

- Using upstream hostname for request \(\#294\) [\#301](https://github.com/fabiolb/fabio/pull/301) ([mitchelldavis](https://github.com/mitchelldavis))

## [v1.4.4](https://github.com/fabiolb/fabio/tree/v1.4.4) (2017-05-08)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.4.3...v1.4.4)

**Implemented enhancements:**

- Add service name to access log fields [\#278](https://github.com/fabiolb/fabio/issues/278)

**Fixed bugs:**

- Fabio does not advertise http/1.1 on TLS connections [\#289](https://github.com/fabiolb/fabio/issues/289)
- fabio does not start with multiple listen sockets [\#279](https://github.com/fabiolb/fabio/issues/279)
- Websocket not working with HTTPS Upstream [\#271](https://github.com/fabiolb/fabio/issues/271)

**Closed issues:**

- Reload configuration without restarting fabio by SIGHUP or by flag. [\#286](https://github.com/fabiolb/fabio/issues/286)
- chunked Transfer-Encoding [\#284](https://github.com/fabiolb/fabio/issues/284)
- How to know what opts are supported in a route / consul tag? [\#270](https://github.com/fabiolb/fabio/issues/270)
- Question: Support for Consul v0.7.3 Node tags [\#252](https://github.com/fabiolb/fabio/issues/252)

## [v1.4.3](https://github.com/fabiolb/fabio/tree/v1.4.3) (2017-04-24)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.4.2...v1.4.3)

**Fixed bugs:**

- Access log cannot be disabled [\#269](https://github.com/fabiolb/fabio/issues/269)

**Closed issues:**

- Can fabio proxy by hostname? [\#267](https://github.com/fabiolb/fabio/issues/267)
- Issues with Haproxy on passthrough mode [\#266](https://github.com/fabiolb/fabio/issues/266)
- How to configure HTTPS upstream manually with tlsskipverify [\#260](https://github.com/fabiolb/fabio/issues/260)
- HTTPS upstream added as HTTP [\#259](https://github.com/fabiolb/fabio/issues/259)

**Merged pull requests:**

- Add support for TLSSkipVerify for https consul fabio check [\#268](https://github.com/fabiolb/fabio/pull/268) ([Ginja](https://github.com/Ginja))

## [v1.4.2](https://github.com/fabiolb/fabio/tree/v1.4.2) (2017-04-10)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.4.1...v1.4.2)

**Implemented enhancements:**

- Add HTTPS upstream support [\#181](https://github.com/fabiolb/fabio/issues/181)

**Closed issues:**

- Find the route across the machine, but no response [\#256](https://github.com/fabiolb/fabio/issues/256)

**Merged pull requests:**

- Allow UI/API to be served over https [\#258](https://github.com/fabiolb/fabio/pull/258) ([tmessi](https://github.com/tmessi))
- Add https upstream support [\#257](https://github.com/fabiolb/fabio/pull/257) ([tmessi](https://github.com/tmessi))

## [v1.4.1](https://github.com/fabiolb/fabio/tree/v1.4.1) (2017-04-04)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.4...v1.4.1)

**Implemented enhancements:**

- Add generic TCP proxying support [\#179](https://github.com/fabiolb/fabio/issues/179)
- Add tests and timeouts to TCP+SNI proxy [\#178](https://github.com/fabiolb/fabio/issues/178)

**Closed issues:**

- Is there any option to enable HSTS [\#254](https://github.com/fabiolb/fabio/issues/254)

## [v1.4](https://github.com/fabiolb/fabio/tree/v1.4) (2017-03-25)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.4rc1...v1.4)

## [v1.4rc1](https://github.com/fabiolb/fabio/tree/v1.4rc1) (2017-03-23)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.4beta2...v1.4rc1)

## [v1.4beta2](https://github.com/fabiolb/fabio/tree/v1.4beta2) (2017-03-23)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.4beta1...v1.4beta2)

## [v1.4beta1](https://github.com/fabiolb/fabio/tree/v1.4beta1) (2017-03-23)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.3.8...v1.4beta1)

**Implemented enhancements:**

- Start listener after routing table is initialized [\#248](https://github.com/fabiolb/fabio/issues/248)
- Support glob host matching [\#163](https://github.com/fabiolb/fabio/issues/163)
- Refactor urlprefix tags [\#111](https://github.com/fabiolb/fabio/issues/111)
- TCP proxying support [\#1](https://github.com/fabiolb/fabio/issues/1)

**Closed issues:**

- feature idea: fabio can be configured to only serve consul services with certain tags [\#245](https://github.com/fabiolb/fabio/issues/245)
- How does services get in to router table of fabio [\#237](https://github.com/fabiolb/fabio/issues/237)

## [v1.3.8](https://github.com/fabiolb/fabio/tree/v1.3.8) (2017-02-14)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.3.7...v1.3.8)

**Implemented enhancements:**

- Retry registry during startup [\#240](https://github.com/fabiolb/fabio/issues/240)
- Make route update logging format configurable [\#238](https://github.com/fabiolb/fabio/issues/238)
- Support absolute URLs [\#219](https://github.com/fabiolb/fabio/issues/219)

**Fixed bugs:**

- requests and notfound metric missing [\#218](https://github.com/fabiolb/fabio/issues/218)
- fabio 1.3.6 UI displays host and path as 'undefined' in the routes page [\#217](https://github.com/fabiolb/fabio/issues/217)

**Closed issues:**

- https support [\#241](https://github.com/fabiolb/fabio/issues/241)
- Fabio - setup details [\#235](https://github.com/fabiolb/fabio/issues/235)
- Not able to connect to fabio UI ... I wonder if I miss any specifics ?. [\#234](https://github.com/fabiolb/fabio/issues/234)
- Error in Fabio setup on container where consul-agent \(client\) is installed [\#233](https://github.com/fabiolb/fabio/issues/233)
- Fabio Connecting error to local consul-agent \(client\) [\#232](https://github.com/fabiolb/fabio/issues/232)
- Load balancing between multiple service cluster nodes [\#231](https://github.com/fabiolb/fabio/issues/231)
- Specify Consul service name in Fabio config [\#230](https://github.com/fabiolb/fabio/issues/230)
- caching [\#228](https://github.com/fabiolb/fabio/issues/228)
- Links in docs to the Traffic Shaping page are dead [\#222](https://github.com/fabiolb/fabio/issues/222)
- Overrides API and GUI save KV Store as wrong name [\#220](https://github.com/fabiolb/fabio/issues/220)

## [v1.3.7](https://github.com/fabiolb/fabio/tree/v1.3.7) (2017-01-19)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.3.6...v1.3.7)

**Implemented enhancements:**

- Support deleting routes by tag [\#201](https://github.com/fabiolb/fabio/issues/201)

**Fixed bugs:**

- Fabio does not serve http2 with go \>= 1.7 [\#215](https://github.com/fabiolb/fabio/issues/215)
- Bad statsd mean metric format [\#207](https://github.com/fabiolb/fabio/issues/207)

**Closed issues:**

- Fabio is not able to pick service from consul and not able to update routing table. [\#210](https://github.com/fabiolb/fabio/issues/210)

## [v1.3.6](https://github.com/fabiolb/fabio/tree/v1.3.6) (2017-01-17)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.3.5...v1.3.6)

**Implemented enhancements:**

- Refactor config loader tests [\#199](https://github.com/fabiolb/fabio/issues/199)
- Routing by path [\#164](https://github.com/fabiolb/fabio/issues/164)
- Strip prefix in the forwarded request [\#44](https://github.com/fabiolb/fabio/issues/44)

**Fixed bugs:**

- runtime error: integer divide by zero [\#186](https://github.com/fabiolb/fabio/issues/186)

**Closed issues:**

- fabio proxy for consul not work, log show no route [\#212](https://github.com/fabiolb/fabio/issues/212)
- Consul registration won't disable [\#209](https://github.com/fabiolb/fabio/issues/209)
- Fabio hangs for 30+ seconds for 204 response [\#206](https://github.com/fabiolb/fabio/issues/206)
- Fabio running using Nomad system scheduler breaks Docker.  [\#192](https://github.com/fabiolb/fabio/issues/192)

## [v1.3.5](https://github.com/fabiolb/fabio/tree/v1.3.5) (2016-11-30)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.3.4...v1.3.5)

**Implemented enhancements:**

- fabio --version switch should work just like -v [\#197](https://github.com/fabiolb/fabio/issues/197)
- Remove proxy.header.tls header from inbound request [\#194](https://github.com/fabiolb/fabio/issues/194)
- Support transparent response body compression [\#119](https://github.com/fabiolb/fabio/issues/119)

**Fixed bugs:**

- missing 'cs' in map [\#189](https://github.com/fabiolb/fabio/issues/189)
- WebSockets not working with IE10 - header casing. [\#183](https://github.com/fabiolb/fabio/issues/183)
- Vault CA Certificate [\#182](https://github.com/fabiolb/fabio/issues/182)

**Closed issues:**

- Logs request [\#188](https://github.com/fabiolb/fabio/issues/188)
- Is this the expecting behavior of Fabio with paths? [\#187](https://github.com/fabiolb/fabio/issues/187)
- TCP+SNI support on the same port as HTTPS  [\#169](https://github.com/fabiolb/fabio/issues/169)

## [v1.3.4](https://github.com/fabiolb/fabio/tree/v1.3.4) (2016-10-28)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.3.3...v1.3.4)

## [v1.3.3](https://github.com/fabiolb/fabio/tree/v1.3.3) (2016-10-12)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.3.2...v1.3.3)

**Implemented enhancements:**

- Provide linux/arm and linux/arm64 binaries [\#161](https://github.com/fabiolb/fabio/issues/161)
- Metrics Prefix with templates [\#160](https://github.com/fabiolb/fabio/pull/160) ([md2k](https://github.com/md2k))

**Fixed bugs:**

- TCP+SNI proxy does not work with PROXY protocol [\#177](https://github.com/fabiolb/fabio/issues/177)
- Consul cert store URL with token not parsed correctly [\#172](https://github.com/fabiolb/fabio/issues/172)
- Panic on invalid response [\#159](https://github.com/fabiolb/fabio/issues/159)

**Closed issues:**

- can not see new application added to the same fabio instance [\#176](https://github.com/fabiolb/fabio/issues/176)
- Ridiculous lack for docker documentation [\#175](https://github.com/fabiolb/fabio/issues/175)
- OT: logo for the eBay organization [\#158](https://github.com/fabiolb/fabio/issues/158)

**Merged pull requests:**

- Use Go's net.JoinHostPort which will auto-detect ipv6 [\#167](https://github.com/fabiolb/fabio/pull/167) ([jovandeginste](https://github.com/jovandeginste))

## [v1.3.2](https://github.com/fabiolb/fabio/tree/v1.3.2) (2016-09-11)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.3.1...v1.3.2)

**Fixed bugs:**

- ParseListen may set the wrong protocol [\#157](https://github.com/fabiolb/fabio/issues/157)

## [v1.3.1](https://github.com/fabiolb/fabio/tree/v1.3.1) (2016-09-09)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.3...v1.3.1)

## [v1.3](https://github.com/fabiolb/fabio/tree/v1.3) (2016-09-09)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.2.1...v1.3)

**Implemented enhancements:**

- Add support for Circonus metrics [\#151](https://github.com/fabiolb/fabio/issues/151)
- Support multiple metrics libraries [\#147](https://github.com/fabiolb/fabio/issues/147)
- Is there a way to prevent SSL requests falling back to an unrelated cert? [\#138](https://github.com/fabiolb/fabio/issues/138)
- Vault token should not require 'root' or 'sudo' privileges [\#134](https://github.com/fabiolb/fabio/issues/134)
- Extended metrics [\#125](https://github.com/fabiolb/fabio/issues/125)

**Fixed bugs:**

- fabio fails to start with "\[FATAL\] 1.2. missing 'cs' in cs" [\#146](https://github.com/fabiolb/fabio/issues/146)

**Closed issues:**

- fabio g-rpc [\#156](https://github.com/fabiolb/fabio/issues/156)
- Routing based on Accept Header [\#155](https://github.com/fabiolb/fabio/issues/155)
- not all command-line options seem to do anything [\#152](https://github.com/fabiolb/fabio/issues/152)

## [v1.2.1](https://github.com/fabiolb/fabio/tree/v1.2.1) (2016-08-25)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.2...v1.2.1)

**Implemented enhancements:**

- Server-sent events support [\#129](https://github.com/fabiolb/fabio/issues/129)
- access logging [\#80](https://github.com/fabiolb/fabio/issues/80)
- Support configuration via command line arguments [\#79](https://github.com/fabiolb/fabio/issues/79)
- Support statsd [\#73](https://github.com/fabiolb/fabio/issues/73)
- SSL Certs from Vault [\#70](https://github.com/fabiolb/fabio/issues/70)
- Refactor listener config [\#28](https://github.com/fabiolb/fabio/issues/28)
- Add/remove certificates using API [\#27](https://github.com/fabiolb/fabio/issues/27)

**Fixed bugs:**

- Always deregister from Consul [\#136](https://github.com/fabiolb/fabio/issues/136)

**Closed issues:**

- HA access to the management interface on instances [\#145](https://github.com/fabiolb/fabio/issues/145)
- Fabio is not adding route, but health check is passing [\#142](https://github.com/fabiolb/fabio/issues/142)
- Wrong Destination IP [\#140](https://github.com/fabiolb/fabio/issues/140)
- Having trouble recognizing routes from consul [\#137](https://github.com/fabiolb/fabio/issues/137)

**Merged pull requests:**

- Improve error message on missing trailing slash [\#143](https://github.com/fabiolb/fabio/pull/143) ([juliangamble](https://github.com/juliangamble))
- added statsd support [\#139](https://github.com/fabiolb/fabio/pull/139) ([jshaw86](https://github.com/jshaw86))

## [v1.2](https://github.com/fabiolb/fabio/tree/v1.2) (2016-07-16)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.2rc4...v1.2)

**Fixed bugs:**

- fabio 1.2rc3 panics with -v [\#128](https://github.com/fabiolb/fabio/issues/128)

## [v1.2rc4](https://github.com/fabiolb/fabio/tree/v1.2rc4) (2016-07-13)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.2rc3...v1.2rc4)

## [v1.2rc3](https://github.com/fabiolb/fabio/tree/v1.2rc3) (2016-07-12)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.1.6...v1.2rc3)

## [v1.1.6](https://github.com/fabiolb/fabio/tree/v1.1.6) (2016-07-12)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.2rc2...v1.1.6)

**Implemented enhancements:**

- TLS handshake error: failed to verify client's certificate [\#108](https://github.com/fabiolb/fabio/issues/108)

**Fixed bugs:**

- X-Forwarded-Port should use local port [\#122](https://github.com/fabiolb/fabio/issues/122)

**Closed issues:**

- Path problem [\#124](https://github.com/fabiolb/fabio/issues/124)

## [v1.2rc2](https://github.com/fabiolb/fabio/tree/v1.2rc2) (2016-06-23)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.1.5...v1.2rc2)

## [v1.1.5](https://github.com/fabiolb/fabio/tree/v1.1.5) (2016-06-23)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.2rc1...v1.1.5)

**Implemented enhancements:**

- Allow routes to a service in warning status [\#117](https://github.com/fabiolb/fabio/pull/117) ([erikvanoosten](https://github.com/erikvanoosten))

**Closed issues:**

- Fabio hangs for 30+ seconds for 204 response [\#120](https://github.com/fabiolb/fabio/issues/120)

## [v1.2rc1](https://github.com/fabiolb/fabio/tree/v1.2rc1) (2016-06-15)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.1.4...v1.2rc1)

## [v1.1.4](https://github.com/fabiolb/fabio/tree/v1.1.4) (2016-06-15)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.1.3...v1.1.4)

**Implemented enhancements:**

- Custom status code when no route found [\#107](https://github.com/fabiolb/fabio/issues/107)
- Keep fabio registered in consul [\#100](https://github.com/fabiolb/fabio/issues/100)
- Disable fabio health check in consul [\#99](https://github.com/fabiolb/fabio/issues/99)
- Support PROXY protocol [\#97](https://github.com/fabiolb/fabio/issues/97)

**Closed issues:**

- fabio should expose a /health endpoint  [\#112](https://github.com/fabiolb/fabio/issues/112)
- Go 1.5 issue [\#109](https://github.com/fabiolb/fabio/issues/109)

## [v1.1.3](https://github.com/fabiolb/fabio/tree/v1.1.3) (2016-05-19)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.1.3rc2...v1.1.3)

**Implemented enhancements:**

- Keep sort order in UI stable [\#104](https://github.com/fabiolb/fabio/issues/104)
- Trim whitespace around tag [\#103](https://github.com/fabiolb/fabio/issues/103)
- SNI support? [\#85](https://github.com/fabiolb/fabio/issues/85)

## [v1.1.3rc2](https://github.com/fabiolb/fabio/tree/v1.1.3rc2) (2016-05-14)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.1.3rc1...v1.1.3rc2)

**Implemented enhancements:**

- Add glob path matching \(an alternative to default prefix matching\) [\#93](https://github.com/fabiolb/fabio/pull/93) ([dkong](https://github.com/dkong))

## [v1.1.3rc1](https://github.com/fabiolb/fabio/tree/v1.1.3rc1) (2016-05-09)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.1.2...v1.1.3rc1)

**Implemented enhancements:**

- Improve forward headers [\#98](https://github.com/fabiolb/fabio/issues/98)
- Allow tags for fabio service registration [\#96](https://github.com/fabiolb/fabio/issues/96)
- Expand experimental HTTP API [\#95](https://github.com/fabiolb/fabio/issues/95)
- Drop default port from request [\#90](https://github.com/fabiolb/fabio/issues/90)
- Use Address instead of ServiceAddress? [\#88](https://github.com/fabiolb/fabio/issues/88)
- Expand ${DC} to consul datacenter [\#55](https://github.com/fabiolb/fabio/issues/55)

**Closed issues:**

- proxy handler error channel bug? [\#92](https://github.com/fabiolb/fabio/issues/92)

## [v1.1.2](https://github.com/fabiolb/fabio/tree/v1.1.2) (2016-04-27)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.1.1...v1.1.2)

**Fixed bugs:**

- Deleted routes hide visible routes [\#57](https://github.com/fabiolb/fabio/issues/57)

**Closed issues:**

- Recommended way to bind multiple fabio instances to public IP for HA [\#89](https://github.com/fabiolb/fabio/issues/89)
- Windows support [\#86](https://github.com/fabiolb/fabio/issues/86)
- How to load balance '/'? [\#83](https://github.com/fabiolb/fabio/issues/83)
- register websockets with consul tags [\#82](https://github.com/fabiolb/fabio/issues/82)
- fabio does not respect registry\_consul\_register\_ip from ENV [\#77](https://github.com/fabiolb/fabio/issues/77)
- Not deregistering when consul health status fails  [\#71](https://github.com/fabiolb/fabio/issues/71)
- question: configure through environment variables? [\#68](https://github.com/fabiolb/fabio/issues/68)
- support middleware\(OWIN\) to execute some code before recirection [\#64](https://github.com/fabiolb/fabio/issues/64)

**Merged pull requests:**

- \#77 fix documentaion [\#78](https://github.com/fabiolb/fabio/pull/78) ([sielaq](https://github.com/sielaq))
- Expose the docker ports in Dockerfile [\#76](https://github.com/fabiolb/fabio/pull/76) ([smancke](https://github.com/smancke))
- Overworked header handling. [\#74](https://github.com/fabiolb/fabio/pull/74) ([smancke](https://github.com/smancke))
- Broken link corrected. [\#65](https://github.com/fabiolb/fabio/pull/65) ([jest](https://github.com/jest))

## [v1.1.1](https://github.com/fabiolb/fabio/tree/v1.1.1) (2016-02-22)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.1...v1.1.1)

**Merged pull requests:**

- Fix use of local ip in consul service registration [\#58](https://github.com/fabiolb/fabio/pull/58) ([jeanblanchard](https://github.com/jeanblanchard))

## [v1.1](https://github.com/fabiolb/fabio/tree/v1.1) (2016-02-18)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.1rc1...v1.1)

**Implemented enhancements:**

- Make read and write timeout configurable [\#53](https://github.com/fabiolb/fabio/issues/53)

## [v1.1rc1](https://github.com/fabiolb/fabio/tree/v1.1rc1) (2016-02-15)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.0.9...v1.1rc1)

## [v1.0.9](https://github.com/fabiolb/fabio/tree/v1.0.9) (2016-02-15)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.0.8...v1.0.9)

**Implemented enhancements:**

- Allow configuration of serviceip used during consul registration [\#48](https://github.com/fabiolb/fabio/issues/48)
- Allow configuration via env vars [\#43](https://github.com/fabiolb/fabio/issues/43)
- Cleanup metrics for deleted routes [\#41](https://github.com/fabiolb/fabio/issues/41)
- HTTP2 support with latest Go [\#32](https://github.com/fabiolb/fabio/issues/32)
- Support additional backends [\#12](https://github.com/fabiolb/fabio/issues/12)

**Fixed bugs:**

- Include services with check ids other than 'service:\*' [\#29](https://github.com/fabiolb/fabio/issues/29)

**Closed issues:**

- Move dependencies to vendor path [\#47](https://github.com/fabiolb/fabio/issues/47)
- Add support for Consul ACL token to demo server [\#37](https://github.com/fabiolb/fabio/issues/37)

## [v1.0.8](https://github.com/fabiolb/fabio/tree/v1.0.8) (2016-01-14)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.0.7...v1.0.8)

**Implemented enhancements:**

- Consul ACL Token [\#36](https://github.com/fabiolb/fabio/issues/36)

**Fixed bugs:**

- Detect when consul agent is down [\#26](https://github.com/fabiolb/fabio/issues/26)
- fabio route not removed after consul deregister [\#22](https://github.com/fabiolb/fabio/issues/22)

**Closed issues:**

- Session persistence [\#33](https://github.com/fabiolb/fabio/issues/33)
- Build fails on master/last release tag [\#31](https://github.com/fabiolb/fabio/issues/31)
- Documentation: make build before running ./fabio [\#24](https://github.com/fabiolb/fabio/issues/24)

**Merged pull requests:**

- \[registry\] fallback to given local IP address [\#30](https://github.com/fabiolb/fabio/pull/30) ([doublerebel](https://github.com/doublerebel))

## [v1.0.7](https://github.com/fabiolb/fabio/tree/v1.0.7) (2015-12-13)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.0.6...v1.0.7)

**Fixed bugs:**

- routes not removed when passing empty string [\#23](https://github.com/fabiolb/fabio/issues/23)

**Closed issues:**

- server demo: Consul health check fails [\#21](https://github.com/fabiolb/fabio/issues/21)
- Demo \(shebang, documentation\) [\#20](https://github.com/fabiolb/fabio/issues/20)
- \(Docker\) Error initializing backend. [\#19](https://github.com/fabiolb/fabio/issues/19)

## [v1.0.6](https://github.com/fabiolb/fabio/tree/v1.0.6) (2015-12-01)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.0.5...v1.0.6)

**Implemented enhancements:**

- Filter routing table not on tags [\#16](https://github.com/fabiolb/fabio/issues/16)
- Support websockets [\#9](https://github.com/fabiolb/fabio/issues/9)

**Fixed bugs:**

- Traffic shaping does not match on service name [\#15](https://github.com/fabiolb/fabio/issues/15)

**Closed issues:**

- Manage manual overrides via UI [\#18](https://github.com/fabiolb/fabio/issues/18)

**Merged pull requests:**

- README: fix typos [\#14](https://github.com/fabiolb/fabio/pull/14) ([ceh](https://github.com/ceh))

## [v1.0.5](https://github.com/fabiolb/fabio/tree/v1.0.5) (2015-11-11)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.0.4...v1.0.5)

**Implemented enhancements:**

- Support Forwarded and X-Forwarded-For headers [\#10](https://github.com/fabiolb/fabio/issues/10)

**Merged pull requests:**

- fix vet warning [\#13](https://github.com/fabiolb/fabio/pull/13) ([juliendsv](https://github.com/juliendsv))

## [v1.0.4](https://github.com/fabiolb/fabio/tree/v1.0.4) (2015-11-03)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.0.3...v1.0.4)

**Implemented enhancements:**

- Support SSL/TLS client cert authentication [\#8](https://github.com/fabiolb/fabio/issues/8)

**Closed issues:**

- List among Consul community tools [\#6](https://github.com/fabiolb/fabio/issues/6)

**Merged pull requests:**

- Fixes broken fragment identifier link [\#11](https://github.com/fabiolb/fabio/pull/11) ([budnik](https://github.com/budnik))

## [v1.0.3](https://github.com/fabiolb/fabio/tree/v1.0.3) (2015-10-26)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.0.2...v1.0.3)

**Merged pull requests:**

- Correcting a typo [\#5](https://github.com/fabiolb/fabio/pull/5) ([mdevreugd](https://github.com/mdevreugd))

## [v1.0.2](https://github.com/fabiolb/fabio/tree/v1.0.2) (2015-10-23)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.0.1...v1.0.2)

**Merged pull requests:**

- Honor consul.url and consul.addr from config file [\#3](https://github.com/fabiolb/fabio/pull/3) ([jeinwag](https://github.com/jeinwag))

## [v1.0.1](https://github.com/fabiolb/fabio/tree/v1.0.1) (2015-10-21)

[Full Changelog](https://github.com/fabiolb/fabio/compare/v1.0.0...v1.0.1)



\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
