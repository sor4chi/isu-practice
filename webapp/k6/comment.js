import http from "k6/http";
import { getRandomAccount, url } from "./config.js";
import { check } from "k6";
import { parseHTML } from "k6/html";

export default function () {
  const [ACCOUNT_NAME, PASSWORD] = getRandomAccount();
  const loginRes = http.post(url("/login"), {
    account_name: ACCOUNT_NAME,
    password: PASSWORD,
  });
  check(loginRes, {
    "is status 200": (r) => r.status === 200,
  });

  const res = http.get(url(`/@${ACCOUNT_NAME}`));

  const doc = parseHTML(res.body);

  const token = doc.find("input[name=csrf_token]").first().attr("value");
  const post_id = doc.find("input[name=post_id]").first().attr("value");

  const comment_res = http.post(url("/comment"), {
    post_id: post_id,
    csrf_token: token,
    comment: "Hello, k6!",
  });
  check(comment_res, {
    "is status 200": (r) => r.status === 200,
  });
}
