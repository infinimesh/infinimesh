import {Component} from '@angular/core';

@Component({
  selector: 'ngx-footer',
  styleUrls: ['./footer.component.scss'],
  template: `
    <span>©2019 — <strong>infinimesh, inc</strong> - source code at <a
      href="https://www.github.com/infinimesh/infinimesh" target="_new"><strong
      style="color: white;">GitHub</strong></a></span>
  `,
})
export class FooterComponent {
}
