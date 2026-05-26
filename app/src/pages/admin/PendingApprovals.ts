import { LitElement, html, css } from 'lit';
import { customElement, state } from 'lit/decorators.js';

@customElement('gx-pending-approvals-dashboard')
export class PendingApprovalsDashboard extends LitElement {
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
      background: linear-gradient(135deg, #10b981 0%, #059669 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      margin: 0 0 8px 0;
    }
    .subtitle {
      color: #9ca3af;
      margin: 0;
    }
    .grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
      gap: 24px;
    }
    .card {
      background: rgba(31, 41, 55, 0.5);
      backdrop-filter: blur(12px);
      border: 1px solid rgba(255, 255, 255, 0.05);
      border-radius: 12px;
      padding: 20px;
      display: flex;
      flex-direction: column;
      justify-content: space-between;
    }
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 16px;
    }
    .badge {
      display: inline-flex;
      align-items: center;
      padding: 4px 8px;
      border-radius: 4px;
      font-size: 0.75rem;
      font-weight: 600;
      text-transform: uppercase;
    }
    .badge-pending {
      background: rgba(245, 158, 11, 0.1);
      color: #f59e0b;
      border: 1px solid rgba(245, 158, 11, 0.2);
    }
    .policy-type {
      font-size: 1.1rem;
      font-weight: 600;
      margin: 0 0 4px 0;
    }
    .details {
      font-size: 0.875rem;
      color: #9ca3af;
      margin-bottom: 20px;
      background: rgba(0, 0, 0, 0.2);
      padding: 12px;
      border-radius: 6px;
    }
    .actions {
      display: flex;
      gap: 12px;
    }
    button {
      flex: 1;
      padding: 10px;
      border-radius: 6px;
      font-weight: 600;
      cursor: pointer;
      border: none;
      transition: background 0.2s;
    }
    .btn-approve {
      background: #10b981;
      color: white;
    }
    .btn-approve:hover {
      background: #059669;
    }
    .btn-reject {
      background: #ef4444;
      color: white;
    }
    .btn-reject:hover {
      background: #dc2626;
    }
  `;

  @state() private requests = [
    {
      id: 'req1',
      user: 'sales1@dibbits.ca',
      branch: 'Trenton',
      policy: 'MIN_MARGIN',
      details: 'Line item 2 Belgard Paver quoted at 8.5% margin (Minimum 12.0%)',
      status: 'PENDING'
    },
    {
      id: 'req2',
      user: 'sales2@dibbits.ca',
      branch: 'Kingston',
      policy: 'CREDIT_LIMIT',
      details: 'Order total $14,500 exceeds customer credit limit by $4,500',
      status: 'PENDING'
    }
  ];

  render() {
    return html`
      <div class="header">
        <h1 class="title">Pending Approvals</h1>
        <p class="subtitle">Review and authorize staff permission overrides for margin and credit limit blocks</p>
      </div>

      <div class="grid">
        ${this.requests.map(
          (req) => html`
            <div class="card">
              <div>
                <div class="card-header">
                  <div>
                    <h3 class="policy-type">${req.policy.replace('_', ' ')}</h3>
                    <div style="font-size: 0.8rem; color: #6b7280;">Branch: ${req.branch} · By: ${req.user}</div>
                  </div>
                  <span class="badge badge-pending">${req.status}</span>
                </div>
                <div class="details">${req.details}</div>
              </div>
              <div class="actions">
                <button class="btn-approve" @click=${() => this.decide(req.id, 'APPROVED')}>Approve</button>
                <button class="btn-reject" @click=${() => this.decide(req.id, 'REJECTED')}>Reject</button>
              </div>
            </div>
          `
        )}
      </div>
    `;
  }

  private decide(id: string, decision: 'APPROVED' | 'REJECTED') {
    console.log(`Override request ${id} decided as ${decision}`);
    // Will be wired to POST /api/v1/staff/approvals/{id}/decide
  }
}
