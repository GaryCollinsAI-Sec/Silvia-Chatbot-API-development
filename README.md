<h1>Silvia Chatbot API Development</h1>

<p>
  <strong>Repository:</strong> Silvia-Chatbot-API-development
</p>

<p>
  A security-focused REST API for <strong>Silvia</strong>, the virtual assistant
  for Silver Dragons Academy. The project demonstrates backend API development,
  secure request handling, chatbot security boundaries, middleware security,
  and automated security testing using Go.
</p>

<hr>

<h2>Project Overview</h2>

<p>
  Silvia is designed to answer questions about publicly approved Silver Dragons
  Academy information while maintaining clear security boundaries around
  internal instructions, private information, administrative systems, and
  unsupported information.
</p>

<p>
  The API was developed using a <strong>security-first approach</strong>.
  Security controls were implemented alongside the API functionality rather
  than being treated as a separate step after development.
</p>

<h2>Technology Stack</h2>

<ul>
  <li><strong>Go</strong> — Backend API</li>
  <li><strong>Chi</strong> — HTTP router</li>
  <li><strong>JSON</strong> — API request and response format</li>
  <li><strong>Go testing package</strong> — Automated testing</li>
  <li><strong>React + TypeScript + Vite</strong> — Frontend client</li>
</ul>

<h2>API Architecture</h2>

<pre>
Frontend
   │
   │ POST /api/chat
   ▼
Go API
   │
   ├── Security Headers
   │
   ├── CORS
   │
   ├── Rate Limiting
   │
   └── Chat Handler
          │
          ├── HTTP method validation
          ├── Content-Type validation
          ├── Request body size limit
          ├── Strict JSON decoding
          ├── Message validation
          └── Message length validation
                 │
                 ▼
            Silvia Chatbot
                 │
                 ├── Instruction disclosure protection
                 ├── Prompt injection protection
                 ├── Private information protection
                 ├── Fabrication protection
                 └── Approved knowledge responses
</pre>

<h2>API Endpoints</h2>

<table>
  <thead>
    <tr>
      <th>Method</th>
      <th>Endpoint</th>
      <th>Purpose</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>GET</td>
      <td><code>/api/health</code></td>
      <td>API health check</td>
    </tr>
    <tr>
      <td>POST</td>
      <td><code>/api/chat</code></td>
      <td>Process Silvia chatbot messages</td>
    </tr>
  </tbody>
</table>

<h2>Security Controls Implemented</h2>

<h3>1. Strict JSON Request Validation</h3>

<p>
  The chat endpoint uses Go's JSON decoder with unknown fields disabled.
  Requests containing unexpected fields are rejected instead of silently
  accepting additional data.
</p>

<p>Example rejected request:</p>

<pre>
{
  "message": "What is Taekwondo?",
  "admin": true
}
</pre>

<p>
  This request is rejected with <code>400 Bad Request</code>.
</p>

<h3>2. Request Body Size Limit</h3>

<p>
  The API limits incoming request bodies to <strong>4096 bytes</strong> using
  <code>http.MaxBytesReader</code>.
</p>

<p>
  This helps prevent unnecessarily large request bodies from reaching the
  application layer.
</p>

<h3>3. Message Length Validation</h3>

<p>
  Individual chatbot messages are limited to <strong>1000 characters</strong>.
  Messages exceeding the limit are rejected with <code>400 Bad Request</code>.
</p>

<h3>4. Empty Message Validation</h3>

<p>
  Messages are trimmed and validated before being passed to the chatbot.
  Empty or whitespace-only messages are rejected.
</p>

<h3>5. HTTP Method Restrictions</h3>

<p>
  The chat endpoint accepts <code>POST</code> requests only.
  Unsupported methods are rejected with <code>405 Method Not Allowed</code>.
</p>

<h3>6. Content-Type Validation</h3>

<p>
  The chat endpoint requires:
</p>

<pre>
Content-Type: application/json
</pre>

<p>
  Requests using an unsupported content type are rejected with
  <code>415 Unsupported Media Type</code>.
</p>

<h3>7. CORS Protection</h3>

<p>
  CORS is configured to allow the approved frontend origin:
</p>

<pre>
http://localhost:5173
</pre>

<p>
  Tests verify that approved origins receive the appropriate CORS headers while
  unapproved origins do not receive authorization through CORS.
</p>

<h3>8. Rate Limiting</h3>

<p>
  The <code>/api/chat</code> endpoint is protected by a rate limiter allowing
  <strong>20 requests per minute per client</strong>.
</p>

<p>
  The rate limiter is applied specifically to the chat route rather than
  globally to every endpoint.
</p>

<p>
  Tests verify:
</p>

<ul>
  <li>Requests within the limit are allowed.</li>
  <li>Requests exceeding the limit are rejected.</li>
  <li>Different clients are tracked separately.</li>
  <li>The rate-limit window resets.</li>
  <li>Invalid client addresses are rejected safely.</li>
</ul>

<h3>9. Security Headers</h3>

<p>
  The API currently applies the following security headers:
</p>

<table>
  <thead>
    <tr>
      <th>Header</th>
      <th>Value</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>X-Content-Type-Options</code></td>
      <td><code>nosniff</code></td>
    </tr>
    <tr>
      <td><code>X-Frame-Options</code></td>
      <td><code>DENY</code></td>
    </tr>
    <tr>
      <td><code>Referrer-Policy</code></td>
      <td><code>no-referrer</code></td>
    </tr>
  </tbody>
</table>

<h2>Chatbot Security Boundaries</h2>

<p>
  Silvia does not simply match questions against academy information.
  Security boundary checks are performed before normal knowledge matching.
</p>

<h3>Instruction Disclosure Protection</h3>

<p>
  Silvia rejects requests attempting to obtain system prompts, internal
  instructions, hidden prompts, or other internal configuration.
</p>

<p>Examples tested include:</p>

<ul>
  <li>Requesting system instructions</li>
  <li>Requesting hidden instructions</li>
  <li>Requesting internal prompts</li>
  <li>Attempting to reveal the system prompt</li>
</ul>

<h3>Prompt Injection Protection</h3>

<p>
  Silvia detects common attempts to override or bypass her instructions.
</p>

<p>
  Security checks include patterns such as requests to ignore previous
  instructions, override instructions, bypass instructions, or act without
  instructions.
</p>

<h3>Private Information Protection</h3>

<p>
  Silvia does not provide access to private academy information,
  administrative systems, databases, credentials, student records, or member
  information.
</p>

<p>
  A general question such as <code>"What is a database?"</code> is not treated
  as a database-access request. This distinction is explicitly tested to help
  reduce false positives.
</p>

<h3>Fabrication Protection</h3>

<p>
  Silvia rejects requests to invent or fabricate academy information.
  This is particularly important for information such as pricing, schedules,
  policies, or other information that has not been approved for the chatbot.
</p>

<h2>Testing Strategy</h2>

<p>
  Testing is divided into multiple layers so that individual security
  components and their integration can be evaluated separately.
</p>

<h3>Chatbot Unit Tests</h3>

<p>
  Tests verify Silvia's knowledge responses and security boundaries, including:
</p>

<ul>
  <li>Instruction disclosure attempts</li>
  <li>Prompt injection attempts</li>
  <li>Private information requests</li>
  <li>Database-access requests</li>
  <li>Fabrication requests</li>
  <li>Approved Taekwondo questions</li>
  <li>Benefits questions</li>
  <li>Pricing questions</li>
  <li>Unknown questions</li>
</ul>

<h3>API Handler Tests</h3>

<p>
  Handler-level tests verify:
</p>

<ul>
  <li>Valid requests</li>
  <li>Empty messages</li>
  <li>Oversized messages</li>
  <li>Malformed JSON</li>
  <li>Unknown JSON fields</li>
  <li>Invalid content types</li>
  <li>Unsupported HTTP methods</li>
  <li>Instruction disclosure requests</li>
</ul>

<h3>Middleware Tests</h3>

<p>
  Middleware tests verify:
</p>

<ul>
  <li>CORS behavior</li>
  <li>Rate limiting</li>
  <li>Client isolation</li>
  <li>Rate-limit resets</li>
  <li>Invalid client addresses</li>
  <li>Security headers</li>
</ul>

<h3>Router Integration Test</h3>

<p>
  A router-level integration test was added to verify the actual middleware
  chain used by the application.
</p>

<pre>
SecurityHeaders
      ↓
CORS
      ↓
RateLimiter
      ↓
POST /api/chat
      ↓
chatHandler
      ↓
Silvia
</pre>

<p>
  The integration test verifies that a legitimate frontend-origin request can
  successfully reach Silvia while the expected security headers and CORS
  controls are applied.
</p>

<h2>Current Test Status</h2>

<p>
  The complete Go test suite currently passes:
</p>

<pre>
go test ./... -v -count=1
</pre>

<p>
  Current result:
</p>

<pre>
PASS
</pre>

<p>
  The root API package, chatbot package, middleware package, and router
  integration tests all pass successfully.
</p>

<h2>Security-First Development Approach</h2>

<p>
  This project follows an incremental security engineering workflow:
</p>

<ol>
  <li>Define the API functionality.</li>
  <li>Identify security boundaries.</li>
  <li>Implement validation and defensive controls.</li>
  <li>Write automated security tests.</li>
  <li>Test individual components.</li>
  <li>Test the integrated router.</li>
  <li>Run the complete regression suite.</li>
  <li>Document the security baseline before expanding functionality.</li>
</ol>

<p>
  The goal is to demonstrate that security controls are considered throughout
  the development lifecycle rather than added only after the application has
  been completed.
</p>

<h2>Current Project Status</h2>

<p>
  <strong>API Security Baseline: Complete</strong>
</p>

<ul>
  <li>REST API foundation implemented</li>
  <li>Silvia chatbot service implemented</li>
  <li>Request validation implemented</li>
  <li>CORS protection implemented</li>
  <li>Rate limiting implemented</li>
  <li>Security headers implemented</li>
  <li>Chatbot security boundaries implemented</li>
  <li>Unit tests implemented</li>
  <li>Middleware tests implemented</li>
  <li>Router integration testing implemented</li>
  <li>Full regression suite passing</li>
</ul>

<hr>

<p>
  <strong>Next development phase:</strong> expand the API carefully while
  preserving the existing security tests as a regression baseline.
</p>
