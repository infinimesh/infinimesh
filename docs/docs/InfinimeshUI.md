# Namespaces Management

Here you can find all namespaces you have rights(at least READ) to.

## First look

![Tabs](Images/namespaces/table.jpg?raw=true)

As you can see whole page, except [**Create namespace**](#create-namespace) button, is just a table.

By the first sight it consists of 3 columns:

  1. Namespace Name - **Mark 5**
  2. ID - **Mark 6**
  3. Actions - **Marks 7, 8. 9, 10**

What's more complicated, that's these rows can have various states:

1. Default **Mark 1**
2. Editable **Mark 2**
3. Expanded **Mark 3**

## Default Mode

In this mode you can just enter [**Edit mode**](#edit-mode)(**Mark 7**) and Delete(**Mark 8**) Namespace

## Edit Mode

In this mode you can rename namespace by editing name in textbox(**Mark 11**) and saving it by clicking Save button(**Mark 9**). If you don't want to apply changes tp the name, click Cancel button(**Mark 10**)

## Expanded Mode(access rights editor)

By clicking on Expand Icon(**Mark 4**) or anywhere on the row, except buttons, you can expand access rights editor.
Which is yet another table with columns:

1. User - **Mark 12**
2. Access Level - **Mark 13**
3. Actions **Marks 14 and 15**

Here you can view, add and delete access rights.

## Create Namespace

On top of the page, you can see **Create namespace** button, by clicking on it, you would get a drawer open.
You don't need any special data to create a namespace: just name.



---------------



# Device Registry Management

Here you can see all devices in the current namespace.

Device Registry is also console "home" page.

## First Look

As you sign in to Console, first you would see this

![Tabs](Images/device-registry/initial.jpg?raw=true)

1. Navigation: back and forward
2. Current Namespace Selector. Since user default namespace is not made for devices, you won't be able to see and create any devices here.
3. Device Registry - This Page
4. See [Accounts](https://github.com/slntopp/infinimesh-frontend-doc/wiki/Accounts-Management-Page)
5. See [Namespaces](https://github.com/slntopp/infinimesh-frontend-doc/wiki/Namespaces-Management-Page)
6. See [Account Management](https://github.com/slntopp/infinimesh-frontend-doc/wiki/Current-User-Management)

## Actual Devices Management

![Tabs](Images/device-registry/selected.jpg?raw=true)

Here you can see a bunch of devices. These are "TestFlight" namespace devices.
At this page you can already perform some actions with this devices.

> Clicking on any of these cards will lead you to [**Device Management Page**](https://github.com/slntopp/infinimesh-frontend-doc/wiki/Device-Management-Page)

### Context menu

You can invoke context menu by right-clicking any device card.

![Tabs](Images/device-registry/context-menu.jpg?raw=true)

### Select

As you can see at [**Mark 1**](#device-registry-management-page) device card is hovered, this means it's selected(So device under [**Mark 2**](#device-registry-management-page) is unselected).
You would need to select devices to perform "batch" actions, such as enabling/disabling multiple devices([**Mark 3**](#device-registry-management-page)).

You can select device by clicking Select option in context menu(see [Context-menu](#context-menu)) or performing `Win + Click` or `Cmd + Click` depending on your OS.

Same for deselecting.

You can select or deselect all devices in group or in registry by clicking on **Select** / **Deselect All** button.

### Search

You can serch through registry using search-bar on top of this page([**Mark 5**](#device-registry-management-page)).

You can see prefix with selector here, possible options are:

* Everywhere - Search through all of the fields described below (default) | Key: all
* Names - Search through names | Key: name
* IDs - Search through hex-ids | Key: id
* Tags -  Filter devices containing given tag | Key: tags
* Namespace - Yet useless | Key: namespace

![Tabs](Images/device-registry/search.jpg?raw=true)

As you could see, every search mode has **key**.
It's used to switch search mode by typing. For example, if you would type id:0xf into search-box, search mode will be automaticaly switched to IDs and you'll filter device containing 0xf in ID.

### Group by tags

In order to find devices faster, you could use **group by tags** functionality.
Just toogle **group by tags** switch([**Mark 6**](#device-registry-management-page)) and you'll get your device grouped:

![Tabs](Images/device-registry/grouped.jpg?raw=true)

* Click on **Whole Registry** switch(**Mark 1**) to stop grouping by tags.

* Click on tag name to expand

* Click **Select All** button(**Mark 2**) to select all devices in current tag group.

> You can also enter group by tag mode and focus on particular tag by clicking on it in **Tags** row inside device card

### Create new Device

If current namespace is not the user default one, you can create a new device by clicking on **Add** button([**Mark 4**](#actual-devices-management))

This show you **Device Creation Drawer**:

![Tabs](Images/device-registry/create-device-drawer.jpg?raw=true)

> **Mark 3**: Specify some tags to make it able to search, sort and group devices easier. (See: [Group by tags](#group-by-tags))

You would need to upload **unique** certificate in order to subscribe to your device MQTT messages.
You can do it either by clicking on **Mark 2** to upload `.crt` file or switching to `Paste` mode via **Mark 1** to paste your certificate from clipboard.


---------------


# Device Management

Here you can manage particular device.

## First look

As you click on any device you'll get to this page.

![Tabs](Images/device-registry/base.jpg?raw=true)

**Mark 1** - Refresh device data button.

**Mark 2** - Bulb color shows if device is enabled(green) or disabled(red), acts same at [Device Registry Page](https://github.com/slntopp/infinimesh-frontend-doc/wiki/Devices-Registry-Management-Page)

**Mark 3** - Device name

**Mark 4** - Device ID

**Mark 5** - Enter Edit Mode

## Edit Mode

After clicking on **Edit** button(**Mark 5**), you'll be able to edit device name and tags:

![Tabs](Images/device-registry/edit-mode.jpg?raw=true)


## State Card

After scrolling little bit down, you can see the Device State Card. It has two columns: **Reported** and **Desired** state:

![Tabs](Images/device-registry/state-base.jpg?raw=true)

**Mark 1 - Reported** state is the state received from the device.
Here you can see a last report timestamp and "version" - order number(**Marks 3 and 2**)

Same for **Desired**(**Mark 4**) and **Marks 5 and 6** for desired state version and timestamp.

By clicking on **Edit** button(**Mark 7**) - you enter **Desired** state edit mode(JSON editor - **Mark 1** below) - this is the data to be sent to the device.

![Tabs](Images/device-registry/state-edit-mode.jpg?raw=true)



---------------



# Accounts Management

Here you can find all user accounts you have rights(at least READ) to.

## First look

![Tabs](Images/accounts/table.jpg?raw=true)

As you can see whole page, except [**Create Account**](#create-account) button, is just a table.
Is consists of 5 columns:

 2. Username
 3. User ID (shown only if your current user is admin of the namespace)
 4. Type (User, Admin or Root)
 5. State - enabled(green) or disabled(red)
 6. Actions Button
    By clicking on it you would get a Context Menu, where you can [**Reset Password**](https://github.com/slntopp/infinimesh-frontend-doc/wiki/Miscellaneous#reset-password), toogle(enable/disable) and delete account

## Create Account

On top of the page, you can see [**Create Account**](#create-account) button, by clicking on it, you would get a drawer open:

![Tabs](Images/accounts/create-account-drawer.jpg?raw=true)




---------------



# Current User Account Management

Yes, it't not even a Page yet. Just a couple button on Sider.

So by hovering your pointer on Sider, you will see menu with all available pages, and your username in the end.

It's expandable by clicking:

![Tabs](Images/current-user/sider.jpg.jpg?raw=true)
![Tabs](Images/current-user/sider-expanded.jpg?raw=true)

After expanding for now you can:

1. [**Reset Password**](https://github.com/slntopp/infinimesh-frontend-doc/wiki/Miscellaneous#reset-password)
2. Log out - which will bring you back to login page.



---------------



# Miscellaneous

Here you can read about little tweaks or components being used across console.

## Reset Password

Modal window with usual form:

![Tabs](Images/miscellaneous/reset-password.jpg?raw=true)

## Themes Selector

You can find Themes Selector at the Console Footer

![Tabs](Images/miscellaneous/themes-selector.jpg?raw=true)

infinimesh Console currently has only three common color schemes:

1. Blue-White -- default
    ![Tabs](Images/miscellaneous/default.jpg?raw=true)
2. Dark Blue(Night)
    ![Tabs](Images/miscellaneous/night.jpg?raw=true)
3. Black and White
    ![Tabs](Images/miscellaneous/black-and-white.jpg?raw=true)


