import http from 'http';
import fs from 'fs';
import path from 'path';

// 测试配置
const config = {
  url: 'http://123.56.99.154/',
  concurrencyLevels: [50, 100, 150, 200, 250, 300],
  duration: 60000, // 1分钟
  rampUpTime: 10000, // 10秒
  resultsDir: './tests/test_results'
};

// 确保结果目录存在
if (!fs.existsSync(config.resultsDir)) {
  fs.mkdirSync(config.resultsDir, { recursive: true });
}

// 解析URL
const url = new URL(config.url);
const options = {
  hostname: url.hostname,
  port: url.port || 80,
  path: url.pathname || '/',
  method: 'GET',
  headers: {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36',
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8',
    'Accept-Language': 'zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2'
  }
};

// 测试结果
class TestResult {
  constructor(concurrency) {
    this.concurrency = concurrency;
    this.requests = 0;
    this.successes = 0;
    this.failures = 0;
    this.totalResponseTime = 0;
    this.minResponseTime = Infinity;
    this.maxResponseTime = 0;
    this.startTime = Date.now();
    this.endTime = null;
  }

  addResult(success, responseTime) {
    this.requests++;
    if (success) {
      this.successes++;
      this.totalResponseTime += responseTime;
      this.minResponseTime = Math.min(this.minResponseTime, responseTime);
      this.maxResponseTime = Math.max(this.maxResponseTime, responseTime);
    } else {
      this.failures++;
    }
  }

  complete() {
    this.endTime = Date.now();
    this.duration = this.endTime - this.startTime;
    this.successRate = (this.successes / this.requests) * 100;
    this.avgResponseTime = this.successes > 0 ? this.totalResponseTime / this.successes : 0;
    this.rps = this.requests / (this.duration / 1000);
  }

  toJSON() {
    return {
      concurrency: this.concurrency,
      requests: this.requests,
      successes: this.successes,
      failures: this.failures,
      successRate: this.successRate.toFixed(2),
      avgResponseTime: this.avgResponseTime.toFixed(2),
      minResponseTime: this.minResponseTime,
      maxResponseTime: this.maxResponseTime,
      duration: this.duration,
      rps: this.rps.toFixed(2),
      startTime: new Date(this.startTime).toISOString(),
      endTime: new Date(this.endTime).toISOString()
    };
  }

  print() {
    console.log(`\n=== 并发 ${this.concurrency} 用户测试结果 ===`);
    console.log(`总请求数: ${this.requests}`);
    console.log(`成功数: ${this.successes}`);
    console.log(`失败数: ${this.failures}`);
    console.log(`成功率: ${this.successRate.toFixed(2)}%`);
    console.log(`平均响应时间: ${this.avgResponseTime.toFixed(2)}ms`);
    console.log(`最小响应时间: ${this.minResponseTime}ms`);
    console.log(`最大响应时间: ${this.maxResponseTime}ms`);
    console.log(`测试持续时间: ${(this.duration / 1000).toFixed(2)}s`);
    console.log(`每秒请求数(RPS): ${this.rps.toFixed(2)}`);
  }
}

// 发送单个请求
function sendRequest(result) {
  return new Promise((resolve) => {
    const startTime = Date.now();
    const req = http.request(options, (res) => {
      const endTime = Date.now();
      const responseTime = endTime - startTime;
      const success = res.statusCode >= 200 && res.statusCode < 300;
      result.addResult(success, responseTime);
      resolve();
    });

    req.on('error', (error) => {
      const endTime = Date.now();
      result.addResult(false, endTime - startTime);
      resolve();
    });

    req.end();
  });
}

// 执行单个并发级别的测试
async function runTest(concurrency) {
  console.log(`\n开始测试: ${concurrency} 并发用户`);
  console.log(`Ramp-up 时间: ${config.rampUpTime / 1000}s`);
  console.log(`测试持续时间: ${config.duration / 1000}s`);

  const result = new TestResult(concurrency);
  const startTime = Date.now();
  const endTime = startTime + config.duration + config.rampUpTime;

  // 逐步增加并发用户
  let currentConcurrency = 0;
  const rampUpStep = config.rampUpTime / concurrency;

  while (Date.now() < endTime) {
    // 增加并发用户
    if (currentConcurrency < concurrency && Date.now() < startTime + config.rampUpTime) {
      currentConcurrency++;
      await new Promise(resolve => setTimeout(resolve, rampUpStep));
    }

    // 发送并发请求
    const promises = [];
    for (let i = 0; i < currentConcurrency; i++) {
      promises.push(sendRequest(result));
    }

    await Promise.all(promises);

    // 短暂休息，避免过度占用系统资源
    await new Promise(resolve => setTimeout(resolve, 50));
  }

  result.complete();
  result.print();

  // 保存结果
  const resultFile = path.join(config.resultsDir, `test_results_${concurrency}.json`);
  fs.writeFileSync(resultFile, JSON.stringify(result.toJSON(), null, 2));
  console.log(`\n测试结果已保存到: ${resultFile}`);

  return result;
}

// 主测试函数
async function main() {
  console.log('=== 服务器承载能力测试 ===');
  console.log(`测试目标: ${config.url}`);
  console.log(`测试并发级别: ${config.concurrencyLevels.join(', ')}`);
  console.log(`每个级别的测试持续时间: ${config.duration / 1000}s`);
  console.log(`每个级别的Ramp-up时间: ${config.rampUpTime / 1000}s`);
  console.log('============================\n');

  const allResults = [];

  for (const concurrency of config.concurrencyLevels) {
    const result = await runTest(concurrency);
    allResults.push(result.toJSON());

    // 测试之间休息一下
    console.log('\n----------------------------');
    console.log('准备下一个测试...');
    console.log('----------------------------\n');
    await new Promise(resolve => setTimeout(resolve, 2000));
  }

  // 保存汇总结果
  const summaryFile = path.join(config.resultsDir, 'test_summary.json');
  fs.writeFileSync(summaryFile, JSON.stringify(allResults, null, 2));
  console.log('\n=== 测试完成 ===');
  console.log(`汇总结果已保存到: ${summaryFile}`);
  console.log('\n测试结果分析:');
  console.log('1. 查看 tests/test_results 目录下的详细结果文件');
  console.log('2. 分析不同并发级别下的响应时间和错误率');
  console.log('3. 找出性能拐点，确定服务器承载上限');
}

// 运行测试
main().catch(error => {
  console.error('测试过程中出现错误:', error);
});