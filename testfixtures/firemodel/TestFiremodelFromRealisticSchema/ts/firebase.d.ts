import * as FIREBASE from 'firebase';

declare module 'firebase' {
  namespace firestore {
    // Snapshots
    export interface DocumentSnapshot<T = DocumentData> {
      data(options?: SnapshotOptions): D | undefined;
    }
    export interface QueryDocumentSnapshot<T = DocumentData> extends DocumentSnapshot {
      data(options?: SnapshotOptions): T;
    }
    export interface QuerySnapshot<T = DocumentData> {
      readonly docs: QueryDocumentSnapshot<T>[];
      forEach(callback: (result: QueryDocumentSnapshot<T>) => void, thisArg?: any): void;
    }

    // References + Queries
    export interface DocumentReference<T = DocumentData> {
      onSnapshot(observer: {
        next?: (snapshot: DocumentSnapshot<T>) => void;
        error?: (error: FirestoreError) => void;
        complete?: () => void;
      }): () => void;
      onSnapshot(
        options: SnapshotListenOptions,
        observer: {
          next?: (snapshot: DocumentSnapshot<T>) => void;
          error?: (error: Error) => void;
          complete?: () => void;
        },
      ): () => void;
      onSnapshot(
        onNext: (snapshot: DocumentSnapshot<T>) => void,
        onError?: (error: Error) => void,
        onCompletion?: () => void,
      ): () => void;
      onSnapshot(
        options: SnapshotListenOptions,
        onNext: (snapshot: DocumentSnapshot<T>) => void,
        onError?: (error: Error) => void,
        onCompletion?: () => void,
      ): () => void;
    }
    export interface Query<T = DocumentData> {
      onSnapshot(observer: {
        next?: (snapshot: QuerySnapshot<T>) => void;
        error?: (error: Error) => void;
        complete?: () => void;
      }): () => void;
      onSnapshot(
        options: SnapshotListenOptions,
        observer: {
          next?: (snapshot: QuerySnapshot<T>) => void;
          error?: (error: Error) => void;
          complete?: () => void;
        },
      ): () => void;
      onSnapshot(
        onNext: (snapshot: QuerySnapshot<T>) => void,
        onError?: (error: Error) => void,
        onCompletion?: () => void,
      ): () => void;
      onSnapshot(
        options: SnapshotListenOptions,
        onNext: (snapshot: QuerySnapshot<T>) => void,
        onError?: (error: Error) => void,
        onCompletion?: () => void,
      ): () => void;
    }
    export interface CollectionReference<T = DocumentData> extends Query<T> {}
  }
}