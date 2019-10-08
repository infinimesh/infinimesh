import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {NamespaceService} from '../../core/data/namespace.service';
import {Observable} from 'rxjs';
import {DeviceService} from '../../core/data/device.service';

@Component({
  selector: 'device-editor',
  templateUrl: './device-editor.component.html',
  styleUrls: ['./device-editor.component.scss']
})
export class DeviceEditorComponent implements OnInit {

  @Input() deviceId;
  @Output() closed = new EventEmitter<boolean>();

  public deviceFormModel;
  public namespaces$: Observable<any>;
  public new = false;

  constructor(private namespaceService: NamespaceService,
              private deviceService: DeviceService) { }

  ngOnInit() {
    this.namespaces$ = this.namespaceService.getAll();
    if(this.deviceId) {
      this.new = false;
      // todo: implement edit
    } else {
      this.new = true;

      this.deviceFormModel = {
        enabled: true,
        certificate: {
          pem_data: this.exampleCertificate
        }
      }
    }
  }

  save() {
    let tags = [];
    if(this.deviceFormModel.tags && this.deviceFormModel.tags.length > 0) {
      tags = this.deviceFormModel.tags.split(',').map(tag => tag.trim());
    }

    const device = {
      name: this.deviceFormModel.name,
      enabled: this.deviceFormModel.enabled,
      namespace: this.deviceFormModel.namespace.name,
      tags: tags,
      certificate: this.deviceFormModel.certificate
    }
    this.deviceService.create(device).subscribe(result => {
      console.log(result);
      this.closed.emit(true);
    });
  }

  private exampleCertificate = 'Paste your Device Certificate';
}
