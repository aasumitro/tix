import * as React from 'react';

export function Footer() {
  return (
    <footer className="fixed bottom-0 w-full bg-white">
      <div className="border-t">
        <div className="flex h-16 items-center px-4">
          <a href="https://github.com/aasumitro/tix" target="_blank" className="mr-2" rel="noreferrer">
            © {new Date().getFullYear()} - TIX
          </a>
          <div className="ml-auto flex items-center space-x-4">
            <ul className="flex flex-row gap-4">
              <li><a href="https://github.com/aasumitro/tix" target="_blank" rel="noreferrer">Docs</a></li>
              <li><a href="https://github.com/aasumitro/tix" target="_blank" rel="noreferrer">FAQs</a></li>
              <li><a href="https://github.com/aasumitro/tix" target="_blank" rel="noreferrer">Support</a></li>
            </ul>
          </div>
        </div>
      </div>
    </footer>
  )
}