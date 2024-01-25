import { expect, test, describe, vi } from "vitest";
import {
  check_offline,
  check_token_expired,
  access_role_conv,
  access_lvl_conv,
} from "@/utils/access";
import { Level, Role } from "infinimesh-proto/build/es/node/access/access_pb";
import { ConnectError } from "@connectrpc/connect";

const check_offline_store = { offline: vi.fn() };
const check_token_expired_store = { logout: vi.fn() };

describe("check_offline", () => {
  test("check is online", () => {
    const offlineSpy = vi.spyOn(check_offline_store, "offline");

    check_offline(new Error("Another error"), check_offline_store);

    expect(offlineSpy).not.toHaveBeenCalled();
  });

  test("check is offline", () => {
    const offlineSpy = vi.spyOn(check_offline_store, "offline");

    check_offline(new Error("Failed to fetch"), check_offline_store);

    expect(offlineSpy).toHaveBeenCalledOnce();
  });

  test("check undefined", () => {
    const offlineSpy = vi.spyOn(check_offline_store, "offline");

    check_offline(undefined, check_offline_store);

    expect(offlineSpy).not.toHaveBeenCalled();
  });
});

describe("check_token_expired", () => {
  test("check expired session", () => {
    const expiredSpy = vi.spyOn(check_token_expired_store, "logout");

    check_token_expired(
      new ConnectError("Session is expired", 16),
      check_token_expired_store
    );

    expect(expiredSpy).toHaveBeenCalledOnce();
  });

  test("check expired token", () => {
    const expiredSpy = vi.spyOn(check_token_expired_store, "logout");

    check_token_expired(
      new ConnectError("Invalid token format", 2),
      check_token_expired_store
    );

    expect(expiredSpy).toHaveBeenCalledOnce();
  });

  test("check another error", () => {
    const expiredSpy = vi.spyOn(check_token_expired_store, "logout");

    check_token_expired(new ConnectError("Test error", 2), check_token_expired_store);

    expect(expiredSpy).not.toHaveBeenCalled();
  });

  test("check undefined", () => {
    const expiredSpy = vi.spyOn(check_token_expired_store, "logout");

    check_token_expired(undefined, check_token_expired_store);

    expect(expiredSpy).not.toHaveBeenCalled();
  });
});

describe("access_role_conv", () => {
  test("check is owner", () => {
    const item = { access: { role: Role.OWNER } };
    expect(access_role_conv(item)).toBe(1);
  });

  test("check is unset", () => {
    const item = { access: { role: Role.UNSET } };
    expect(access_role_conv(item)).toBe(0);
  });

  test("check is role number", () => {
    const item = { access: { role: 1 } };

    expect(access_role_conv(item)).toBe(1);
  });

  test("check is role wrong number", () => {
    const item = { access: { role: 5 } };

    expect(access_role_conv(item)).toBe(0);
  });

  test("check undefined", () => {
    expect(access_role_conv()).toBe(0);
  });
});

describe("access_lvl_conv", () => {
  test("check is ADMIN", () => {
    const item = { access: { level: Level.ADMIN } };
    expect(access_lvl_conv(item)).toBe(3);
  });

  test("check is NONE", () => {
    const item = { access: { level: Level.NONE } };
    expect(access_lvl_conv(item)).toBe(0);
  });

  test("check is level number", () => {
    const item = { access: { level: 3 } };

    expect(access_lvl_conv(item)).toBe(3);
  });

  test("check is level wrong number", () => {
    const item = { access: { level: 8 } };

    expect(access_lvl_conv(item)).toBe(1);
  });

  test("check undefined", () => {
    expect(access_lvl_conv()).toBe(1);
  });
});
