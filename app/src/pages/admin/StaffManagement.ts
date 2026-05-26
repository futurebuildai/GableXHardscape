import { LitElement, html, css } from 'lit';
import { customElement, state } from 'lit/decorators.js';

@customElement('gx-staff-management')
export class StaffManagement extends LitElement {
  static styles = css`
    :host {
      display: block;
      padding: 24px;
      color: #f3f4f6;
    }
    .header {
      margin-bottom: 32px;
    }
    .title {
      font-size: 2rem;
      font-weight: 700;
      background: linear-gradient(135deg, #fbbf24 0%, #d97706 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      margin: 0 0 8px 0;
    }
    .subtitle {
      color: #9ca3af;
      margin: 0;
    }
    .table-container {
      background: rgba(31, 41, 55, 0.5);
      backdrop-filter: blur(12px);
      border: 1px solid rgba(255, 255, 255, 0.05);
      border-radius: 12px;
      overflow: hidden;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      text-align: left;
    }
    th {
      background: rgba(17, 24, 39, 0.7);
      padding: 16px;
      font-weight: 600;
      color: #9ca3af;
      border-bottom: 1px solid rgba(255, 255, 255, 0.05);
    }
    td {
      padding: 16px;
      border-bottom: 1px solid rgba(255, 255, 255, 0.05);
    }
    tr:hover {
      background: rgba(255, 255, 255, 0.02);
    }
    .role-badge {
      display: inline-flex;
      align-items: center;
      padding: 4px 10px;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 600;
    }
    .role-super {
      background: rgba(245, 158, 11, 0.1);
      color: #f59e0b;
      border: 1px solid rgba(245, 158, 11, 0.2);
    }
    .role-manager {
      background: rgba(59, 130, 246, 0.1);
      color: #3b82f6;
      border: 1px solid rgba(59, 130, 246, 0.2);
    }
    .role-core {
      background: rgba(16, 185, 129, 0.1);
      color: #10b981;
      border: 1px solid rgba(16, 185, 129, 0.2);
    }
    select {
      background: #1f2937;
      border: 1px solid rgba(255, 255, 255, 0.1);
      color: white;
      padding: 6px 12px;
      border-radius: 6px;
      outline: none;
      cursor: pointer;
    }
    select:focus {
      border-color: #f59e0b;
    }
  `;

  @state() private users = [
    { sub: 'auth0|user1', email: 'colton@dibbits.ca', role: 'General Manager', branches: ['Trenton', 'Kingston'] },
    { sub: 'auth0|user2', email: 'amanda@dibbits.ca', role: 'Financial Controller', branches: ['Trenton', 'Kingston'] },
    { sub: 'auth0|user3', email: 'tanya@dibbits.ca', role: 'Sales Manager', branches: ['Trenton'] },
    { sub: 'auth0|user4', email: 'yard1@dibbits.ca', role: 'Yard Team', branches: ['Kingston'] }
  ];

  render() {
    return html`
      <div class="header">
        <h1 class="title">Staff & Roles</h1>
        <p class="subtitle">Assign roles and branch scoping constraints across the 13-role permissions hierarchy</p>
      </div>

      <div class="table-container">
        <table>
          <thead>
            <tr>
              <th>User Sub / Email</th>
              <th>Current Role</th>
              <th>Assigned Branches</th>
              <th>Manage Role</th>
            </tr>
          </thead>
          <tbody>
            ${this.users.map(
              (user) => html`
                <tr>
                  <td>
                    <div style="font-weight: 600;">${user.email}</div>
                    <div style="font-size: 0.75rem; color: #6b7280;">${user.sub}</div>
                  </td>
                  <td>
                    <span class="role-badge ${this.getRoleClass(user.role)}">
                      ${user.role}
                    </span>
                  </td>
                  <td>${user.branches.join(', ')}</td>
                  <td>
                    <select @change=${(e: Event) => this.updateRole(user.sub, (e.target as HTMLSelectElement).value)}>
                      <optgroup label="Admin Tier (Cross-Branch)">
                        <option value="General Manager" ?selected=${user.role === 'General Manager'}>General Manager</option>
                        <option value="Financial Controller" ?selected=${user.role === 'Financial Controller'}>Financial Controller</option>
                      </optgroup>
                      <optgroup label="Manager Tier (Scoped)">
                        <option value="Branch Manager" ?selected=${user.role === 'Branch Manager'}>Branch Manager</option>
                        <option value="Procurement Manager" ?selected=${user.role === 'Procurement Manager'}>Procurement Manager</option>
                        <option value="Sales Manager" ?selected=${user.role === 'Sales Manager'}>Sales Manager</option>
                        <option value="Yard Manager" ?selected=${user.role === 'Yard Manager'}>Yard Manager</option>
                        <option value="Logistics Manager" ?selected=${user.role === 'Logistics Manager'}>Logistics Manager</option>
                        <option value="HR" ?selected=${user.role === 'HR'}>HR</option>
                      </optgroup>
                      <optgroup label="Staff Tier (Scoped)">
                        <option value="Inside Sales" ?selected=${user.role === 'Inside Sales'}>Inside Sales</option>
                        <option value="Outside Sales" ?selected=${user.role === 'Outside Sales'}>Outside Sales</option>
                        <option value="Yard Team" ?selected=${user.role === 'Yard Team'}>Yard Team</option>
                        <option value="Drivers" ?selected=${user.role === 'Drivers'}>Drivers</option>
                        <option value="Payables/Receivables" ?selected=${user.role === 'Payables/Receivables'}>Payables/Receivables</option>
                      </optgroup>
                    </select>
                  </td>
                </tr>
              `
            )}
          </tbody>
        </table>
      </div>
    `;
  }

  private getRoleClass(role: string): string {
    const adminRoles = ['General Manager', 'Financial Controller'];
    const managerRoles = ['Branch Manager', 'Procurement Manager', 'Sales Manager', 'Yard Manager', 'Logistics Manager', 'HR'];
    
    if (adminRoles.includes(role)) return 'role-super';
    if (managerRoles.includes(role)) return 'role-manager';
    return 'role-core';
  }

  private updateRole(sub: string, newRole: string) {
    console.log(`Updating role of ${sub} to ${newRole}`);
    // Will be wired to PUT /api/v1/staff/roles/{sub}
  }
}
