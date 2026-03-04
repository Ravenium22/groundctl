import React from 'react';
import Layout from '@theme/Layout';
import Link from '@docusaurus/Link';

export default function Home() {
  return (
    <Layout title="Home" description="terraform plan for your local developer machine">
      <main style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: 'calc(100vh - 200px)',
        padding: '2rem',
        textAlign: 'center',
      }}>
        <h1 style={{ fontSize: '3rem', marginBottom: '0.5rem' }}>groundctl</h1>
        <p style={{ fontSize: '1.4rem', color: 'var(--ifm-color-emphasis-700)', maxWidth: '600px' }}>
          <code>terraform plan</code> for your local developer machine.
        </p>
        <p style={{ fontSize: '1.1rem', color: 'var(--ifm-color-emphasis-600)', maxWidth: '600px', marginBottom: '2rem' }}>
          Detect how your development environment has drifted from your team's standard — and fix it with one command.
        </p>
        <div style={{ display: 'flex', gap: '1rem' }}>
          <Link
            className="button button--primary button--lg"
            to="/docs/getting-started">
            Get Started
          </Link>
          <Link
            className="button button--secondary button--lg"
            href="https://github.com/Ravenium22/groundctl">
            GitHub
          </Link>
        </div>
        <pre style={{
          marginTop: '3rem',
          textAlign: 'left',
          padding: '1.5rem 2rem',
          borderRadius: '8px',
          fontSize: '0.95rem',
          maxWidth: '500px',
          width: '100%',
        }}>
{`$ ground check

  [ok]  node          22.10.0
  [ok]  python        3.12.1
  [ERR] docker        not found
  [!!]  terraform     version drift

  4 checked  2 ok  1 warning  1 error

$ ground fix
  Installing docker via brew...
  Upgrading terraform via brew...
  All drift resolved.`}
        </pre>
      </main>
    </Layout>
  );
}
