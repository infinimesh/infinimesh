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
export function check_token_expired(err, store) {
  if (err.response && err.response.status == 500) {
    if (err.response.data.message == "Token is expired") {
      store.logout();
    }
  }
}