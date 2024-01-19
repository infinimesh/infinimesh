export const access_levels = {
  "NONE": 0,
  "READ": 1,
  "MGMT": 2,
  "ADMIN": 3,
  "ROOT": 4,
}

export function access_lvl_conv(o) {
  return access_levels[o.access.level || "READ"];
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