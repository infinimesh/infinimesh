// describe("Login Page render tests", () => {
//   setupTest({
//     server: true
//   });

//   it("redirects from / to /login", async () => {
//     const { page } = await get("/");
//     console.log(page);
//   });

//   // it("renders html", () => {
//   //   // ...
//   // });
// });

import { shallowMount } from "@vue/test-utils";
import Login from "../login";

import Vue from "vue";
import Antd from "ant-design-vue";
import vuetify from "vuetify";

describe("Login", () => {
  let wrapper;
  beforeEach(() => {
    Vue.use(Antd);
    Vue.use(vuetify);

    wrapper = shallowMount(Login);
  });

  it("mounts properly", () => {
    expect(wrapper.vm).toBeTruthy();
  });

  it("renders properly", () => {
    expect(wrapper.html()).toContain("<h1>infinimesh</h1>");
  });
});
