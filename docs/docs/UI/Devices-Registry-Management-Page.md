# Device Registry Management

Here you can see all devices in the current namespace.

Device Registry is also console "home" page.

## First Look

As you sign in to Console, first you would see this

![Tabs](Images/device-registry/initial.jpg?raw=true)

1. Navigation: back and forward
2. Current Namespace Selector. Since user default namespace is not made for devices, you won't be able to see and create any devices here.
3. Device Registry - This Page
4. See [Accounts](Accounts-Management-Page.md)
5. See [Namespaces](Namespaces-Management-Page.md)
6. See [Account Management](Current-User-Management.md)

## Actual Devices Management

![Tabs](Images/device-registry/selected.jpg?raw=true)

Here you can see a bunch of devices. These are "TestFlight" namespace devices.
At this page you can already perform some actions with this devices.

> Clicking on any of these cards will lead you to [**Device Management Page**](Device-Management-Page.md)

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