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

  private exampleCertificate = '-----BEGIN CERTIFICATE-----\n' +
    'MIIEljCCAn4CCQD26Bj66/x6ZjANBgkqhkiG9w0BAQsFADANMQswCQYDVQQGEwJk\n' +
    'ZTAeFw0xOTA0MTQyMDQ0MDZaFw0yMDA0MTMyMDQ0MDZaMA0xCzAJBgNVBAYTAmRl\n' +
    'MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAzU4CvLI8+R3La9+uYlGE\n' +
    'xuGiiyyiDeEXOzP8xve+ojKK++3LNWrljo2l5qqpAyFiGhbFOfKDfMe4DvSnAp5C\n' +
    'mzCV0J+6VBLtHYouFixesO5ZL6/aoB3XZwtr4PP5/aQ3y0s3RcdSfVbsEL6TY+Zc\n' +
    'BQEIn7Xyjrbrx6+O8dz0jECfksdH0OKrlxJ/6xOOHFHnAnWX31WbTnEefnUa4U1N\n' +
    '0nA68X875X+SDPhNSvAV3FRL7G2hTDvTnWO1hLZ6OICcMMi6EI+lR1syIfWxpYon\n' +
    'v7HbPEBeS/8R+orEQhVCOi3AmZLMW2ZEDU4/urKA6LI5WRmZI/zdC3A9onp5RfGS\n' +
    'K/+glRA54YkQynhmZ7Gu4VD1Cmn4gsb/AQE+BS1ORHz9n2fTwU8rkLcVIPFIaq0z\n' +
    'yh1b2JvYLVhcg0ajkvNNMAVjD/T1FPL7P3JnDLtOnnQ/z/yS9o2vKo0rydD8OU3B\n' +
    'd4yb7c9fSnvXMB32JwtPwKdvRhpfNWGaaP4uKozrxF3tuksM/7RRmtKINbrgDKns\n' +
    'CM0ppVo1qe7mc6xcY8tJQLxDDncPx8V1fg7NbBIy1TsjX1i/fL/tBUXhGRkDiDEm\n' +
    'vuKRMQORkIpWT89uPGeYd1ATZCOFW6BAePhChfLHSv+MIclIqGRBXyalq7Vks4w8\n' +
    '+ywR68E7MrOhwBvRAMYz4tMCAwEAATANBgkqhkiG9w0BAQsFAAOCAgEAg2HbmNTc\n' +
    'VS9JqJYlyzpzLencRp4cVuGgBIHi+RJitFLsxTt7MlvudOqSy2gutaM7RGNc/xNA\n' +
    '6QaiTlYoPEREGyuT5F+rhSzNn5rdqHcU7pbcwI1FNoKK3FmKcScx4mwq40cd7CDP\n' +
    '7LesH8gDGktvY+LIoZM2pIUJsoJjXoG4blW1fb8fML3GUQSqiEohWnNpw8uYdpre\n' +
    'lQcke1AeTyVXIF7rTVR7N48pc9u2m+U/8I765jiZdXeZS2wphBSCsDfsV/LeKJ97\n' +
    '/M8g7LZMFfIh9hxURJ0kB4mV7Pm9W4RxXji2SNa+iizGk9UmUH0rDCPzyWGrDpT1\n' +
    'RGiX8SybLLnQuKd41YSqxnJBA5/LTvTRYUVBGNj26KSLBulCq64FmR88EP/EVqoC\n' +
    'M6I5NCmeGYIfaaJEB9U/d90vibBfY4p4ir8FHAwdzFXQo7ckTov/uCb70a5mZaqR\n' +
    'rXYOL/n/BPehqFeGeazF+SH8PN4pYZR8kRp+CgCteFqC1Oze8OGJOho+CuwBJgWK\n' +
    'bqL/LfQQyxgDJXfRnKeU3Pm/7syXbb+aVFS4+ENMbM+pLv37sCvS7JoMD8meiCft\n' +
    'qfjuBEVnym5Kf7gTK5awPiAzyHPFMlWb6QldgObCnYWqOg+KuVylFwyCLtzL8YGP\n' +
    'wUC5VPTn2n4kGKLZvBog+iXz4OVkf9xhH4U=\n' +
    '-----END CERTIFICATE-----';
}
