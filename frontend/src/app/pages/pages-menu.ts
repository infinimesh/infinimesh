import {NbMenuItem} from '@nebular/theme';

export const MENU_ITEMS: NbMenuItem[] = [
  /*{
    title: 'Dashboard',
    icon: 'nb-home',
    link: '/pages/dashboard',
    home: true,
  },*/
  {
    title: 'Device Registry',
    icon: 'nb-e-commerce',
    link: '/pages/devices',
    pathMatch: '/devices/.*'
  }
];
