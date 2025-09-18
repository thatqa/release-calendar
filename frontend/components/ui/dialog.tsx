"use client";

import * as React from "react";

type BaseProps = {
    children: React.ReactNode;
    className?: string;
};

export function Dialog({
                           open,
                           onOpenChange,
                           children,
                       }: {
    open: boolean;
    onOpenChange: (v: boolean) => void;
    children: React.ReactNode;
}) {
    if (!open) return null;
    return (
        <div
            className="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
            onClick={() => onOpenChange(false)}
        >
            {/* stop propagation so clicks inside content don't close */}
            <div onClick={(e) => e.stopPropagation()}>{children}</div>
        </div>
    );
}

export function DialogContent({children, className}: BaseProps) {
    return (
        <div
            className={
                "bg-white rounded-lg p-6 w-full max-w-lg shadow-lg " +
                (className || "")
            }
        >
            {children}
        </div>
    );
}

export function DialogHeader({children, className}: BaseProps) {
    return <div className={"mb-2 " + (className || "")}>{children}</div>;
}

export function DialogTitle({children, className}: BaseProps) {
    return (
        <h2 className={"text-lg font-semibold " + (className || "")}>{children}</h2>
    );
}

export function DialogDescription({children, className}: BaseProps) {
    return (
        <p className={"text-sm text-neutral-600 " + (className || "")}>
            {children}
        </p>
    );
}
