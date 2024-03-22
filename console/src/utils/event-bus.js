export class EventBus {
  #eventsMap;

  constructor() {
    this.#eventsMap = new Map();
  }

  publish(eventName, ...args) {
    const callbackObject = this.#eventsMap.get(eventName);

    if (!callbackObject) return;
    for (let id of Object.getOwnPropertySymbols(callbackObject)) {
      callbackObject[id](...args);
    }
  }

  subscribe(eventName, callback) {
    if (!this.#eventsMap.get(eventName)) {
      this.#eventsMap.set(eventName, {});
    }

    const id = Symbol(eventName);

    this.#eventsMap.get(eventName)[id] = callback;
    const unsubscribe = () => {
      delete this.#eventsMap.get(eventName)[id];

      if (
        Object.getOwnPropertySymbols(this.#eventsMap.get(eventName)).length ===
        0
      ) {
        this.#eventsMap.delete(eventName);
      }
    };

    return { unsubscribe };
  }
}
