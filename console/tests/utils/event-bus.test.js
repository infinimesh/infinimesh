import { expect, test, describe, vi } from "vitest";
import { EventBus } from "@/utils/event-bus";

describe("event_bus", () => {
  test("base events", () => {
    const callbacks = { firstCallback: vi.fn(), secondCallback: vi.fn() };

    const firstCallbackSpy = vi.spyOn(callbacks, "firstCallback");
    const secondCallbackSpy = vi.spyOn(callbacks, "secondCallback");

    const eventBus = new EventBus();

    eventBus.subscribe("first", callbacks.firstCallback);
    eventBus.subscribe("second", callbacks.secondCallback);

    eventBus.publish("first", 1);
    eventBus.publish("second", 2);

    expect(firstCallbackSpy).toHaveBeenCalledOnce();
    expect(firstCallbackSpy).toHaveBeenCalledWith(1);
    expect(secondCallbackSpy).toHaveBeenCalledOnce();
    expect(secondCallbackSpy).toHaveBeenCalledWith(2);
  });

  test("many callbacks on one event", () => {
    const callbacks = {
      firstCallback: vi.fn(),
      secondCallback: vi.fn(),
      thirdCallback: vi.fn(),
    };

    const firstCallbackSpy = vi.spyOn(callbacks, "firstCallback");
    const secondCallbackSpy = vi.spyOn(callbacks, "secondCallback");
    const thirdCallbackSpy = vi.spyOn(callbacks, "thirdCallback");

    const eventBus = new EventBus();

    eventBus.subscribe("event", callbacks.firstCallback);
    eventBus.subscribe("event", callbacks.secondCallback);
    eventBus.subscribe("event", callbacks.thirdCallback);

    eventBus.publish("event", "first");
    eventBus.publish("event");
    eventBus.publish("event", "last");

    expect(firstCallbackSpy).toHaveBeenCalledTimes(3);
    expect(secondCallbackSpy).toHaveBeenCalledTimes(3);
    expect(thirdCallbackSpy).toHaveBeenCalledTimes(3);

    expect(thirdCallbackSpy).toHaveBeenCalledWith("first");
    expect(thirdCallbackSpy).toHaveBeenLastCalledWith("last");
  });

  test("unsubscribe", () => {
    const callbacks = { firstCallback: vi.fn() };

    const firstCallbackSpy = vi.spyOn(callbacks, "firstCallback");

    const eventBus = new EventBus();

    const { unsubscribe } = eventBus.subscribe(
      "first",
      callbacks.firstCallback
    );

    eventBus.publish("first", 1);

    unsubscribe();

    eventBus.publish("first", 2);

    expect(firstCallbackSpy).toHaveBeenCalledOnce();
    expect(firstCallbackSpy).toHaveBeenCalledWith(1);
    expect(firstCallbackSpy).not.toHaveBeenCalledWith(2);
  });

  test("unsubscribe + many events", () => {
    const callbacks = {
      firstCallback: vi.fn(),
      secondCallback: vi.fn(),
      thirdCallback: vi.fn(),
    };

    const firstCallbackSpy = vi.spyOn(callbacks, "firstCallback");
    const secondCallbackSpy = vi.spyOn(callbacks, "secondCallback");
    const thirdCallbackSpy = vi.spyOn(callbacks, "thirdCallback");

    const eventBus = new EventBus();

    const { unsubscribe: unsubscribeFirst } = eventBus.subscribe(
      "event",
      callbacks.firstCallback
    );
    const { unsubscribe: unsubscribeSecond } = eventBus.subscribe(
      "event",
      callbacks.secondCallback
    );
    const { unsubscribe: unsubscribeThird } = eventBus.subscribe(
      "event",
      callbacks.thirdCallback
    );

    eventBus.publish("event", "first");
    unsubscribeFirst();

    eventBus.publish("event");
    unsubscribeSecond();

    eventBus.publish("event", "last");
    unsubscribeThird();

    expect(firstCallbackSpy).toHaveBeenCalledTimes(1);
    expect(secondCallbackSpy).toHaveBeenCalledTimes(2);
    expect(thirdCallbackSpy).toHaveBeenCalledTimes(3);
  });
});
