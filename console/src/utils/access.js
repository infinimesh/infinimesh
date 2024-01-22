import { Role, Level } from "infinimesh-proto/build/es/node/access/access_pb"

export function access_lvl_conv(item) {
  const { level = "READ" } = item.access;

  if (level >= 0 && level <= 4) return level
  return Level[level] ?? 1;
}

export function access_role_conv(item) {
  const { role = "UNSET" } = item.access;

  return Role[role] ?? Role.UNSET;
}

// must be axios error and App Store
// will remove in next time
export function check_token_expired_http(err, store) {
  if (err.response && err.response.status == 500) {
    if (err.response.data.message == "Token is expired") {
      store.logout({
        title: "Signed Out",
        type: "warning",
        description: "Token has expired! Please log in again."
      });
    }
  }

  if (err.response && err.response.status == 401) {
    if (err.response.data.message.includes("Session is expired")) {
      store.logout({
        title: "Signed Out",
        type: "warning",
        description: "Session has expired or has been revoked! Please log in again."
      });
    }
  }
}

export function check_token_expired(err) {
  if (err.code == 2 && err.message.toLowerCase().includes('Invalid token format')) {
    store.logout({
      title: "Signed Out",
      type: "warning",
      description: "Token has expired! Please log in again."
    });
  }

  if (err.code == 16 && err.message === 'Session is expired, revoked or invalid') {
    store.logout({
      title: "Signed Out",
      type: "warning",
      description: "Session has expired or has been revoked! Please log in again."
    });
  }
}

// will remove in next time
export function check_offline_http(err, store) {
  if (!err.isAxiosError) {
    return
  }

  if (err.message == "Network Error") {
    store.offline()
  }
}

export function check_offline(err, store) {
  if (err.message == "Failed to fetch") {
    store.offline()
  }
}