import * as React from 'react';
export function Select({className = '', children, ...props}: React.SelectHTMLAttributes<HTMLSelectElement>) { return <select className={`select ${className}`} {...props}>{children}</select>; }
