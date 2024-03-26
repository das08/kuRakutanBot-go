import { describe, it } from "mocha";
import expect from "expect.js";
import axios from "axios";

axios.defaults.baseURL = "http://localhost:3000";

describe("API", function () {
  it("should pass", function () {
    expect("it works").to.be("it works");
  });
});
