package malio

import (
	"fmt"
	"testing"
)

func TestBracketExtraction(t *testing.T) {
	content := `<script>
  function sendData() {
    $crisp.push(["set", "user:email", "xxx@gmail.com"], ["set", "user:nickname", "xxx"]);
    $crisp.push(["set", "session:data", [[
      ["Class", "1"],
      ["Class_Expire", "2021-03-30 13:43:31"],
      ["Money", "100.49"],
      ["Unused_Traffic", "91.2GB"]
    ]]]);
  }
  sendData();
</script>`
	if key := BracketExtraction(content, "Money"); key != "100.49" {
		fmt.Println(key)
		t.Fail()
	}

	if key := BracketExtraction(content, "user:email"); key != "xxx@gmail.com" {
		fmt.Println(key)
		t.Fail()
	}
}

func TestTodayUsed(t *testing.T) {
	content := `<div class="card-wrap">
                      <div class="card-header">
                        <h4>剩余流量</h4>
                      </div>
                      <div class="card-body">
                        <span class="counter">91.06</span> GB
                      </div>
                      <div class="card-stats">
                        <div class="card-stats-title" style="padding-top: 0;padding-bottom: 4px;">
                          <nav aria-label="breadcrumb">
                            <ol class="breadcrumb">
                              <li class="breadcrumb-item active" aria-current="page">今日已用: 989.31MB<br><br>重置时间: 未购买套餐.</li>
                            </ol>
                          </nav>
                        </div>
                      </div>`
	if mb := TodayUsed(content); mb != "989.31MB" {
		fmt.Println(mb)
		t.Fail()
	}

}

func TestCardDescExtraction(t *testing.T) {
	content := `      </div>
          </div>
                                                              <div class="row">
                <div class="col-lg-3 col-md-3 col-sm-12">
                  <div class="card card-statistic-2">
                    <div class="card-icon shadow-primary bg-primary">
                      <i class="fas fa-crown"></i>
                    </div>
                    <div class="card-wrap">
                      <div class="card-header">
                        <h4>会员时长</h4>
                      </div>
                      <div class="card-body">
                                                <span class="counter">25</span> 天
                                              </div>
                    </div>
                    <div class="card-stats">
                      <div class="card-stats-title" style="padding-top: 0;padding-bottom: 4px;">
                        <nav aria-label="breadcrumb">
                          <ol class="breadcrumb">
                            <li class="breadcrumb-item active" aria-current="page">
                                                              VIP①:
                              
                                                              2021-03-30 过期
                              <br><br>最高带宽: 100M`
	if de := CardDescExtraction(content,"会员时长");de != 25 {
		fmt.Println(de)
		t.Fail()
	}
}
