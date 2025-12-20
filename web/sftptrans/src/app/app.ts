import { Component, OnInit, HostListener } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { FilePanel } from './components/file-panel/file-panel';
import { ApiService, FileEntry, StatusResponse } from './services/api.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, HttpClientModule, FilePanel],
  template: `
    <div class="app-container">
      <header class="header">
        <span class="title">sftptrans</span>
        <span class="connection" *ngIf="status">
          [Connected: {{ status.connection }}]
        </span>
        <button class="exit-btn" (click)="onExit()">Exit</button>
      </header>

      <main class="main-content">
        <app-file-panel
          title="Local Files"
          [files]="localFiles"
          [currentPath]="localPath"
          (pathChange)="onLocalPathChange($event)"
          (fileSelect)="onLocalFileSelect($event)"
          [selectedFile]="selectedLocalFile"
        >
          <div class="actions">
            <button (click)="onUpload()" [disabled]="!selectedLocalFile || selectedLocalFile.isDir">
              Upload Selected →
            </button>
          </div>
        </app-file-panel>

        <app-file-panel
          title="Remote Files"
          [files]="remoteFiles"
          [currentPath]="remotePath"
          (pathChange)="onRemotePathChange($event)"
          (fileSelect)="onRemoteFileSelect($event)"
          [selectedFile]="selectedRemoteFile"
        >
          <div class="actions">
            <button (click)="onDownload()" [disabled]="!selectedRemoteFile || selectedRemoteFile.isDir">
              ← Download Selected
            </button>
            <button (click)="onNewFolder()">New Folder</button>
            <button (click)="onRename()" [disabled]="!selectedRemoteFile">Rename</button>
            <button (click)="onDelete()" [disabled]="!selectedRemoteFile" class="danger">Delete</button>
          </div>
        </app-file-panel>
      </main>

      <footer class="status-bar">
        Status: {{ statusMessage }}
      </footer>
    </div>
  `,
  styles: [`
    .app-container {
      display: flex;
      flex-direction: column;
      height: 100vh;
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
      font-size: 14px;
    }
    .header {
      display: flex;
      align-items: center;
      padding: 8px 16px;
      background: #2c3e50;
      color: white;
      gap: 16px;
    }
    .title {
      font-weight: bold;
      font-size: 18px;
    }
    .connection {
      color: #7fdbca;
      font-family: monospace;
    }
    .exit-btn {
      margin-left: auto;
      background: #e74c3c;
      color: white;
      border: none;
      padding: 6px 16px;
      cursor: pointer;
    }
    .exit-btn:hover {
      background: #c0392b;
    }
    .main-content {
      display: flex;
      flex: 1;
      overflow: hidden;
    }
    .actions {
      display: flex;
      gap: 8px;
      padding: 8px;
      background: #ecf0f1;
      flex-wrap: wrap;
    }
    .actions button {
      padding: 6px 12px;
      cursor: pointer;
      background: #3498db;
      color: white;
      border: none;
    }
    .actions button:hover:not(:disabled) {
      background: #2980b9;
    }
    .actions button:disabled {
      background: #bdc3c7;
      cursor: not-allowed;
    }
    .actions button.danger {
      background: #e74c3c;
    }
    .actions button.danger:hover:not(:disabled) {
      background: #c0392b;
    }
    .status-bar {
      padding: 8px 16px;
      background: #34495e;
      color: #ecf0f1;
      font-family: monospace;
    }
  `]
})
export class AppComponent implements OnInit {
  status: StatusResponse | null = null;
  statusMessage = 'Ready';

  localFiles: FileEntry[] = [];
  remoteFiles: FileEntry[] = [];
  localPath = '';
  remotePath = '/';
  selectedLocalFile: FileEntry | null = null;
  selectedRemoteFile: FileEntry | null = null;

  constructor(private api: ApiService) {}

  ngOnInit() {
    this.loadStatus();
    this.loadRemoteFiles('/');
  }

  @HostListener('window:keydown', ['$event'])
  handleKeyDown(event: KeyboardEvent) {
    if (event.ctrlKey && event.key === 'q') {
      event.preventDefault();
      this.onExit();
    }
    if (event.key === 'F5') {
      event.preventDefault();
      this.refresh();
    }
  }

  loadStatus() {
    this.api.getStatus().subscribe({
      next: (res) => {
        this.status = res.data;
        this.localPath = res.data.downloadDir;
        this.loadLocalFiles(this.localPath);
      },
      error: (err) => {
        this.statusMessage = 'Error: ' + err.message;
      }
    });
  }

  loadLocalFiles(path: string) {
    this.api.listLocal(path).subscribe({
      next: (res) => {
        this.localFiles = res.data || [];
        this.localPath = path;
        this.selectedLocalFile = null;
      },
      error: (err) => {
        this.statusMessage = 'Error loading local files: ' + err.error?.error?.message;
      }
    });
  }

  loadRemoteFiles(path: string) {
    this.statusMessage = 'Loading...';
    this.api.listRemote(path).subscribe({
      next: (res) => {
        this.remoteFiles = res.data || [];
        this.remotePath = path;
        this.selectedRemoteFile = null;
        this.statusMessage = 'Ready';
      },
      error: (err) => {
        this.statusMessage = 'Error: ' + (err.error?.error?.message || err.message);
      }
    });
  }

  onLocalPathChange(path: string) {
    this.loadLocalFiles(path);
  }

  onRemotePathChange(path: string) {
    this.loadRemoteFiles(path);
  }

  onLocalFileSelect(file: FileEntry) {
    this.selectedLocalFile = file;
  }

  onRemoteFileSelect(file: FileEntry) {
    this.selectedRemoteFile = file;
  }

  onUpload() {
    if (!this.selectedLocalFile) return;
    this.statusMessage = 'Uploading...';
    this.api.upload(this.selectedLocalFile.path, this.remotePath).subscribe({
      next: () => {
        this.statusMessage = 'Upload complete';
        this.loadRemoteFiles(this.remotePath);
      },
      error: (err) => {
        this.statusMessage = 'Upload failed: ' + err.error?.error?.message;
      }
    });
  }

  onDownload() {
    if (!this.selectedRemoteFile) return;
    this.statusMessage = 'Downloading...';
    this.api.download(this.selectedRemoteFile.path).subscribe({
      next: (res) => {
        this.statusMessage = 'Downloaded to: ' + res.data.localPath;
        this.loadLocalFiles(this.localPath);
      },
      error: (err) => {
        this.statusMessage = 'Download failed: ' + err.error?.error?.message;
      }
    });
  }

  onNewFolder() {
    const name = prompt('Enter folder name:');
    if (!name) return;
    const path = this.remotePath === '/' ? '/' + name : this.remotePath + '/' + name;
    this.api.mkdir(path).subscribe({
      next: () => {
        this.loadRemoteFiles(this.remotePath);
      },
      error: (err) => {
        this.statusMessage = 'Error: ' + err.error?.error?.message;
      }
    });
  }

  onRename() {
    if (!this.selectedRemoteFile) return;
    const newName = prompt('Enter new name:', this.selectedRemoteFile.name);
    if (!newName || newName === this.selectedRemoteFile.name) return;

    const parentPath = this.selectedRemoteFile.path.substring(
      0,
      this.selectedRemoteFile.path.lastIndexOf('/')
    ) || '/';
    const newPath = parentPath === '/' ? '/' + newName : parentPath + '/' + newName;

    this.api.rename(this.selectedRemoteFile.path, newPath).subscribe({
      next: () => {
        this.loadRemoteFiles(this.remotePath);
      },
      error: (err) => {
        this.statusMessage = 'Error: ' + err.error?.error?.message;
      }
    });
  }

  onDelete() {
    if (!this.selectedRemoteFile) return;
    if (!confirm(`Delete "${this.selectedRemoteFile.name}"?`)) return;

    this.api.delete(this.selectedRemoteFile.path).subscribe({
      next: () => {
        this.loadRemoteFiles(this.remotePath);
      },
      error: (err) => {
        this.statusMessage = 'Error: ' + err.error?.error?.message;
      }
    });
  }

  onExit() {
    if (!confirm('Disconnect and exit?')) return;
    this.api.shutdown().subscribe({
      next: () => {
        this.statusMessage = 'Disconnected. You can close this tab.';
      }
    });
  }

  refresh() {
    this.loadLocalFiles(this.localPath);
    this.loadRemoteFiles(this.remotePath);
  }
}
