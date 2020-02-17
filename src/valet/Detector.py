#!/usr/local/bin/python
import argparse
from collections import OrderedDict

import cv2
import dlib
import imutils
import numpy as np
from scipy.spatial import distance as dist

CONSIDER_CLASSES = ["background", "aeroplane", "bicycle", "bird", "boat",
           "bottle", "bus", "car", "cat", "chair", "cow", "diningtable",
           "dog", "horse", "motorbike", "person", "pottedplant", "sheep",
           "sofa", "train", "tvmonitor"]


class CentroidTracker:

    def __init__(self, maxDisappeared=50, maxDistance=50):
        self.nextObjectID = 0
        self.objects = OrderedDict()
        self.disappeared = OrderedDict()
        self.maxDisappeared = maxDisappeared
        self.maxDistance = maxDistance

    def register(self, centroid):
        self.objects[self.nextObjectID] = centroid
        self.disappeared[self.nextObjectID] = 0
        self.nextObjectID += 1

    def deregister(self, objectID):
        del self.objects[objectID]
        del self.disappeared[objectID]

    def update(self, rects):
        if len(rects) == 0:
            for objectID in list(self.disappeared.keys()):
                self.disappeared[objectID] += 1
                if self.disappeared[objectID] > self.maxDisappeared:
                    self.deregister(objectID)
            return self.objects

        inputCentroids = np.zeros((len(rects), 2), dtype="int")

        for (i, (startX, startY, endX, endY)) in enumerate(rects):
            cX = int((startX + endX) / 2.0)
            cY = int((startY + endY) / 2.0)
            inputCentroids[i] = (cX, cY)

        if len(self.objects) == 0:
            for i in range(0, len(inputCentroids)):
                self.register(inputCentroids[i])
        else:
            objectIDs = list(self.objects.keys())
            objectCentroids = list(self.objects.values())
            D = dist.cdist(np.array(objectCentroids), inputCentroids)
            rows = D.min(axis=1).argsort()
            cols = D.argmin(axis=1)[rows]
            usedRows = set()
            usedCols = set()

            for (row, col) in zip(rows, cols):
                if row in usedRows or col in usedCols:
                    continue

                if D[row, col] > self.maxDistance:
                    continue

                objectID = objectIDs[row]
                self.objects[objectID] = inputCentroids[col]
                self.disappeared[objectID] = 0

                usedRows.add(row)
                usedCols.add(col)

            unusedRows = set(range(0, D.shape[0])).difference(usedRows)
            unusedCols = set(range(0, D.shape[1])).difference(usedCols)

            if D.shape[0] >= D.shape[1]:
                for row in unusedRows:
                    objectID = objectIDs[row]
                    self.disappeared[objectID] += 1

                    if self.disappeared[objectID] > self.maxDisappeared:
                        self.deregister(objectID)

            else:
                for col in unusedCols:
                    self.register(inputCentroids[col])

        return self.objects


class TrackableObject:

    def __init__(self, objectID):
        self.objectID = objectID
        self.centroids = []
        self.counted = False

    def appendCentroid(self, centroid):
        self.centroids.append(centroid)


def monitor(files):
    # Maintain a count of cars, motorbikes, or boats seen in these frames
    # although only one stream of video is analyzed, there may be multiple cars
    count = 0
    centroidTracker = CentroidTracker(maxDisappeared=40, maxDistance=50)
    trackableObjects = {}
    confidenceLevel = 0.9
    consider = ['car', 'motorbike', 'boat']
    trainedModel = cv2.dnn.readNetFromCaffe('../../assets/deploy.prototxt', '../../assets/deploy.caffemodel')

    # Iterate over each frame passed to the script
    for file in files:
        frame = cv2.imread(file)
        frame = imutils.resize(frame, width=400)
        (height, width) = frame.shape[:2]
        rgb = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
        rects = []
        trackers = []
        blob = cv2.dnn.blobFromImage(
            cv2.resize(frame, (300, 300)),
            0.007843,
            (300, 300),
            127.5
        )
        trainedModel.setInput(blob)
        detections = trainedModel.forward()

        for i in np.arange(0, detections.shape[2]):
            confidence = detections[0, 0, i, 2]

            if confidence > confidenceLevel:
                idx = int(detections[0, 0, i, 1])

                if CONSIDER_CLASSES[idx] in consider:
                    box = detections[0, 0, i, 3:7] * np.array([width, height, width, height])
                    (startX, startY, endX, endY) = box.astype("int")
                    tracker = dlib.correlation_tracker()
                    rect = dlib.rectangle(startX, startY, endX, endY)
                    tracker.start_track(rgb, rect)

                    trackers.append(tracker)
            else:
                for tracker in trackers:
                    tracker.update(rgb)
                    pos = tracker.get_position()
                    rects.append(
                        (int(pos.left()),
                         int(pos.top()),
                         int(pos.right()),
                         int(pos.bottom()))
                    )

        cv2.line(frame, (width // 2, 0), (width // 2, height), (0, 0, 255), 2)

        objects = centroidTracker.update(rects)

        for (objectID, centroid) in objects.items():
            object = trackableObjects.get(objectID, None)
            if object is None:
                object = TrackableObject(objectID)
            else:
                direction = centroid[0] - np.mean([centroid[0] for centroid in object.centroids])
                object.appendCentroid(centroid)

                if not object.counted:
                    if direction < -20 and centroid[0] < width // 2:
                        count += 1
                        object.counted = True
                    elif direction > 20 and centroid[0] > width // 2:
                        count -= 1
                        object.counted = True

            trackableObjects[objectID] = object

    # This line will only be reached if we did not detect an entry or exit
    return count


def main():
    # Pass filenames to detector via args
    parser = argparse.ArgumentParser()
    parser.add_argument('files', nargs='+')
    args = parser.parse_args()

    # return monitor(args.files)
    # testing only
    print(monitor(args.files))


if __name__ == '__main__':
    main()
