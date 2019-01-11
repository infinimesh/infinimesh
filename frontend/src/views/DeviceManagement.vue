<template>
  <v-container>
    <h1 class="mb-3">Device Management</h1>
    <v-card>
      <v-treeview
        :items="tree"
      ></v-treeview>
    </v-card>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      counter: 0,
      tree: [],
      tree: [
          {
            id: 1,
            name: 'Applications',
            children: [
              { id: 2, name: 'Calendar' },
              { id: 3, name: 'Chrome' },
              { id: 4, name: 'Webstorm' }
            ]
          },
          {
            id: 5,
            name: 'Languages',
            children: [
              { id: 6, name: 'English' },
              { id: 7, name: 'French' },
              { id: 8, name: 'Spannish' }
            ]
          }

        ],
      realTree: `
      {
       "objects": {
         "0x11181": {
           "uid": "0x11181",
           "name": "Johannes' Home",
           "objects": {
             "0x11182": {
               "uid": "0x11182",
               "name": "First Floor",
               "objects": {
                 "0x11183": {
                   "uid": "0x11183",
                   "name": "Living Room",
                   "devices": [
                     {
                       "uid": "0x11184",
                       "name": "PC"
                     }
                   ]
                 }
               }
             },
             "0x11187": {
               "uid": "0x11187",
               "name": "Second Floor"
             }
           },
           "devices": [
             {
               "uid": "0x11188",
               "name": "le lamp"
             }
           ]
         }
       }
      }
      `,
      simpleTree: {
        objects: {
          0x11183: {
            uid: "0x11183",
            name: "Living Room",
            devices: [
              {
                uid: "0x11184",
                name: "PC"
              }
            ]
          }
        }
      },
      easyObj: {
        topLevel: {
          one: 1,
          two: 2
        }
      }
    }
  },
  computed: {
    objectTree() {
      return JSON.parse(this.realTree)
    }
  },
  methods: {
    constructTreeView(obj) {
      console.log("object at beginning of function", obj)
      let newObj = {};
      console.log("newObj at beginning of function", newObj)
      console.log("start construct tree")
      console.log("loop number", this.counter)

      if (this.isObject(obj)) {
        let key;
        console.log("in object part")

        //counter is only there atm to prevent an infinite regression
        this.counter++;
        if (this.counter > 3) {
          return;
        }
        else {
          for (key in obj) {
            obj.children = [];
            console.log("key", key)
            console.log("value for key", obj[key])
            obj.children.push(obj[key]);
            console.log("obj.children after push", obj.children)
            //copy object
            newObj = JSON.parse(JSON.stringify(obj.children))
            console.log("newObj", newObj, "obj", obj)
            obj = {};
            return this.constructTreeView(newObj);
          }
        }
      }
      else if (this.isArray(obj)) {
        console.log("in array part")

        obj.children = [];
        this.counter++;
        if (this.counter > 3) {
          return;
        }
        obj.children = obj;
        newObj = JSON.parse(JSON.stringify(obj.children))
        console.log(obj)
        obj = {};
        return this.constructTreeView(newObj);
      }
      else {
        console.log("done", obj);
      }
    },
    isArray(a) {
    return (!!a) && (a.constructor === Array);
    },
    isObject(a) {
    return (!!a) && (a.constructor === Object);
    }
  },
  mounted() {
    console.log(this.objectTree);
    this.constructTreeView(this.easyObj);
  }
}
</script>

<style>
</style>
