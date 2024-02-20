import { describe, test, expect, vi } from "vitest";
import { mount } from "@vue/test-utils";
import { useAppStore } from "@/store/app";
import { useAccountsStore } from "@/store/accounts";
import { defineComponent, markRaw } from "vue";
import { createPinia, setActivePinia } from "pinia";

const TestComponent = defineComponent({
  template: "<div></div>",
});

const pinia = createPinia();

const mockBar = { start: vi.fn, finish: vi.fn(), error: vi.fn() };

describe("store/accounts", () => {
  beforeEach(() => {
    const wrapper = mount(TestComponent, {
      global: {
        plugins: [pinia],
      },
    });

    setActivePinia(pinia);
  });

  test("sync_me", async () => {
    const appStore = useAppStore();
    const accountsStore = useAccountsStore();

    await accountsStore.sync_me();

    expect(appStore.me.title).not.toBe("");
    expect(appStore.me.uuid).not.toBeNull();
  });

  test("fetchAccounts", async () => {
    const accountsStore = useAccountsStore();

    const promise = accountsStore.fetchAccounts();
    expect(accountsStore.loading).toBe(true);

    await promise;

    expect(accountsStore.accounts).not.toBe({});
    expect(accountsStore.loading).toBe(false);
  });

  test("updateAccount", async () => {
    const startSpy = vi.spyOn(mockBar, "start");
    const finishSpy = vi.spyOn(mockBar, "finish");

    const accountsStore = useAccountsStore();

    const result = await accountsStore.updateAccount(
      { title: "test" },
      mockBar
    );

    expect(!!result).toBe(false);
    expect(startSpy).toHaveBeenCalledOnce();
    expect(finishSpy).toHaveBeenCalledOnce();
  });

  test("createAccount", async () => {
    const startSpy = vi.spyOn(mockBar, "start");
    const finishSpy = vi.spyOn(mockBar, "finish");

    const accountsStore = useAccountsStore();
    const accounts = JSON.parse(JSON.stringify(accountsStore.accounts));

    await accountsStore.createAccount({ title: "test" }, mockBar);

    expect(accountsStore.accounts).not.toBe(accounts);
    expect(startSpy).toHaveBeenCalledOnce();
    expect(finishSpy).toHaveBeenCalledOnce();
  });

  test("toggle", async () => {
    const startSpy = vi.spyOn(mockBar, "start");
    const finishSpy = vi.spyOn(mockBar, "finish");

    const accountsStore = useAccountsStore();
    const accounts = JSON.parse(JSON.stringify(accountsStore.accounts));

    await accountsStore.toggle({ title: "test" }, mockBar);

    expect(accountsStore.accounts).not.toBe(accounts);
    expect(startSpy).toHaveBeenCalledOnce();
    expect(finishSpy).toHaveBeenCalledOnce();
  });

  test("deleteAccount", async () => {
    const startSpy = vi.spyOn(mockBar, "start");
    const finishSpy = vi.spyOn(mockBar, "finish");

    const accountsStore = useAccountsStore();
    const accounts = JSON.parse(JSON.stringify(accountsStore.accounts));

    await accountsStore.deleteAccount("test", mockBar);

    expect(accountsStore.accounts).not.toBe(accounts);
    expect(startSpy).toHaveBeenCalledOnce();
    expect(finishSpy).toHaveBeenCalledOnce();
  });

  test("move", async () => {
    const accountsStore = useAccountsStore();
    const account = Object.values(accountsStore.accounts)[0];
    const testNamespace = { access: "test" };
    await accountsStore.moveAccount(account.uuid, testNamespace);

    expect(accountsStore.accounts[account.uuid].access.namespace).toStrictEqual(
      testNamespace
    );
  });

  test("setCredentials", async () => {
    const startSpy = vi.spyOn(mockBar, "start");
    const finishSpy = vi.spyOn(mockBar, "finish");

    const accountsStore = useAccountsStore();
    accountsStore.accounts=[]
    await accountsStore.fetchAccounts();
    const accounts = JSON.parse(JSON.stringify(accountsStore.accounts));

    await accountsStore.setCredentials("test", "test", mockBar);

    expect(accountsStore.accounts).not.toBe(accounts);
    expect(startSpy).toHaveBeenCalledOnce();
    expect(finishSpy).toHaveBeenCalledOnce();
  });

  test("tokenFor", async () => {
    const accountsStore = useAccountsStore();
    const account = Object.values(accountsStore.accounts)[0];

    const { token } = await accountsStore.tokenFor(account);
    expect(token).toBeTypeOf("string");
  });
});
