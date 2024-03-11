import { mount } from "@vue/test-utils";
import { ref } from "vue";
import DeviceSearch from "@/components/dashboard/device-search.vue";

// Test suite for device-search
describe("DeviceSearch", () => {
  it("renders the component", () => {
    const wrapper = mount(DeviceSearch);

    expect(wrapper.exists()).toBe(true);
  });

  it("renders the placeholder", () => {
    const wrapper = mount(DeviceSearch);
    const el = wrapper.find(".n-base-selection-placeholder__inner");

    // get placeholder attribute
    const placeholder = el.text();
    expect(placeholder).toBe("Filter devices eg. :uuid:abc");
  });

  it("doesn't render the placeholder when there is value", async () => {
    const key = "uuid:04cd0083-4329-45e8-873b-3be20fb130ef";

    const wrapper = mount(DeviceSearch, {
      props: {
        filterTerm: [key],
      },
    });

    const el = wrapper.find(".n-base-selection-placeholder__inner");
    expect(el.exists()).toBe(false);
  });

  it("renders the tags as passed", async () => {
    const key = "uuid:04cd0083-4329-45e8-873b-3be20fb130ef";

    const wrapper = mount(DeviceSearch, {
      props: {
        filterTerm: [key],
      },
    });

    expect(wrapper.find(".n-tag__content").text()).toBe(key);
  });

  it("emits update:value event when tag is removed", async () => {
    const key = "uuid:04cd0083-4329-45e8-873b-3be20fb130ef";

    const value = ref([key]);

    const wrapper = mount(DeviceSearch, {
      props: {
        filterTerm: value,
      },
    });

    await wrapper.find(".n-tag__close").trigger("click");
    expect(wrapper.emitted("update:value")).toBeTruthy();
  });
});
