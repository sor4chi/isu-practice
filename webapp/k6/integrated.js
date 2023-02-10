import initialize from "./initialize.js";
import comment from "./comment.js";
import postimage from "./postimage.js";

export { initialize, comment, postimage };

// VUs == Virtual Users

export const options = {
  scenarios: {
    initialize: {
      executor: "shared-iterations", // 一定量の実行を複数のVUsで共有する
      vus: 1, // 1VUs
      iterations: 1, // 1回の実行
      exec: "initialize", // initialize.jsを実行
      maxDuration: "10s", // 10秒間実行(レギュレーション)
    },
    comment: {
      executor: "constant-vus", //　複数のVUsを並列で実行する
      vus: 4, // 4VUs
      duration: "30s", // 30秒間繰り返し実行
      exec: "comment", // comment.jsを実行
      startTime: "12s", // 12秒後に実行開始
    },
    postimage: {
      executor: "constant-vus", // 複数のVUsを並列で実行する
      vus: 2, // 2VUs
      duration: "30s", // 30秒間繰り返し実行
      exec: "postimage", // postimageを実行
      startTime: "12s", // 12秒後に実行開始
    },
  },
};
