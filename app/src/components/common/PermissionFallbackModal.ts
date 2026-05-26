import { LitElement, html, css } from 'lit';
import { customElement, property, state } from 'lit/decorators.js';

@customElement('gx-permission-fallback-modal')
export class PermissionFallbackModal extends LitElement {
  static styles = css`
    .modal-overlay {
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(0, 0, 0, 0.7);
      backdrop-filter: blur(8px);
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 1000;
    }
    .modal-container {
      background: #1f2937;
      border: 1px solid rgba(255, 255, 255, 0.1);
      border-radius: 16px;
      width: 100%;
      max-width: 480px;
      padding: 32px;
      color: #f3f4f6;
      box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5);
    }
    .modal-header {
      display: flex;
      align-items: center;
      gap: 16px;
      margin-bottom: 24px;
    }
    .icon-container {
      background: rgba(239, 68, 68, 0.1);
      color: #ef4444;
      padding: 12px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    .title {
      font-size: 1.5rem;
      font-weight: 700;
      margin: 0;
    }
    .description {
      color: #9ca3af;
      font-size: 0.95rem;
      line-height: 1.5;
      margin-bottom: 24px;
    }
    .policy-box {
      background: rgba(0, 0, 0, 0.2);
      border: 1px solid rgba(255, 255, 255, 0.05);
      border-radius: 8px;
      padding: 16px;
      margin-bottom: 24px;
    }
    .policy-title {
      font-size: 0.8rem;
      color: #6b7280;
      text-transform: uppercase;
      font-weight: 700;
      margin-bottom: 8px;
    }
    .policy-detail {
      font-size: 1rem;
      font-weight: 600;
      color: #f3f4f6;
    }
    .textarea-container {
      margin-bottom: 24px;
    }
    label {
      display: block;
      font-size: 0.875rem;
      font-weight: 600;
      color: #9ca3af;
      margin-bottom: 8px;
    }
    textarea {
      width: 100%;
      background: rgba(0, 0, 0, 0.2);
      border: 1px solid rgba(255, 255, 255, 0.1);
      border-radius: 8px;
      color: white;
      padding: 12px;
      outline: none;
      resize: none;
      box-sizing: border-box;
    }
    textarea:focus {
      border-color: #ef4444;
    }
    .actions {
      display: flex;
      gap: 12px;
    }
    button {
      flex: 1;
      padding: 12px;
      border-radius: 8px;
      font-weight: 600;
      cursor: pointer;
      border: none;
    }
    .btn-submit {
      background: #ef4444;
      color: white;
    }
    .btn-submit:hover {
      background: #dc2626;
    }
    .btn-cancel {
      background: rgba(255, 255, 255, 0.05);
      color: #9ca3af;
      border: 1px solid rgba(255, 255, 255, 0.1);
    }
    .btn-cancel:hover {
      background: rgba(255, 255, 255, 0.1);
      color: white;
    }
  `;

  @property({ type: Boolean }) open = false;
  @property({ type: String }) policyType = 'MIN_MARGIN';
  @property({ type: String }) details = 'Attempted to quote line item at 9.0% margin, below the 12.0% threshold.';
  @state() private notes = '';
  @state() private submitted = false;

  render() {
    if (!this.open) return html``;

    return html`
      <div class="modal-overlay">
        <div class="modal-container">
          <div class="modal-header">
            <div class="icon-container">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
            </div>
            <h2 class="title">Action Blocked</h2>
          </div>

          <p class="description">
            Your account is restricted from performing this action due to the branch policy below. You can request a manager override to bypass this block.
          </p>

          <div class="policy-box">
            <div class="policy-title">Policy Restriction</div>
            <div class="policy-detail">${this.details}</div>
          </div>

          ${this.submitted
            ? html`
                <div style="background: rgba(16, 185, 129, 0.1); color: #10b981; padding: 16px; border-radius: 8px; text-align: center; font-weight: 600; margin-bottom: 24px;">
                  Approval request submitted successfully.
                </div>
                <button class="btn-cancel" style="width: 100%;" @click=${this.closeModal}>Close</button>
              `
            : html`
                <div class="textarea-container">
                  <label for="notes">Override Justification Notes</label>
                  <textarea id="notes" rows="3" placeholder="Provide details on why you need this exception..." .value=${this.notes} @input=${(e: Event) => this.notes = (e.target as HTMLTextAreaElement).value}></textarea>
                </div>

                <div class="actions">
                  <button class="btn-cancel" @click=${this.closeModal}>Cancel</button>
                  <button class="btn-submit" @click=${this.submitRequest}>Request Override</button>
                </div>
              `}
        </div>
      </div>
    `;
  }

  private closeModal() {
    this.open = false;
    this.submitted = false;
    this.notes = '';
    this.dispatchEvent(new CustomEvent('close'));
  }

  private submitRequest() {
    this.submitted = true;
    console.log(`Submitting override request for ${this.policyType} with notes: ${this.notes}`);
    // Will be wired to POST /api/v1/staff/approvals
  }
}
