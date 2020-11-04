import { shallowMount } from "@vue/test-utils";
import Index from "../index";

describe("Dashboard index", () => {
  let wrapper;
  let $router;
  let push = jest.fn();

  beforeEach(() => {
    $router = {
      push
    };
    wrapper = shallowMount(Index, { mocks: { $router } });
  });

  it("mounts properly", () => {
    expect(wrapper.vm).toBeTruthy();
  });
  it("redirects", () => {
    expect(push).toHaveBeenCalled();
  });
  it("redirects to dashboard-devices", () => {
    expect(JSON.stringify(push.mock.calls[0][0])).toBe(
      JSON.stringify({ name: "dashboard-devices" })
    );
  });
});
