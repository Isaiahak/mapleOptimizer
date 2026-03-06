import numpy as np
from PIL import Image
import os
import torch
import torch.nn as nn
from torch.utils.data import Dataset
from torch.utils.data import DataLoader

from torchvision import transforms

classes = [
"Hero",
"Pally",
"DK",
"FP",
 "IL",
"Bishop",
"Bowmaster",
"Marksman",
"PathFinder",
"NL",
"Shadower",
"DB",
"Bucc",
"Corsair",
"CM",
"DW",
"BW",
"WA",
"NW",
"TB",
"Aran",
"Evan",
"Mercedes",
"Phantom",
"Shade",
"Luminous",
"DS",
"DA",
"BM",
"WH",
"Mechanic",
"Xenon",
"Blaster",
"Hayato",
"Kanna",
"Mihile",
"Kaiser",
"Kain",
"Cadena",
"AB",
"Zero",
"Saitama",
"Kinesis",
"Adele",
"Illium",
"Khali",
"Ark",
"Ren",
"Lara",
"Hoyoung",
"Lynn",
"MX",
"Sia",
]

class TwoClass(nn.Module):
    def __init__(self):
        super().__init__()
        self.layers = nn.ModuleList([
            nn.Conv2d(3,16, kernel_size=3, padding=1),
            nn.ReLU(),
            nn.MaxPool2d(2,2),
            nn.Conv2d(16,32, kernel_size=3, padding=1),
            nn.ReLU(),
            nn.MaxPool2d(2,2),
            nn.Conv2d(32,64, kernel_size=3, padding=1),
            nn.ReLU(),
            nn.MaxPool2d(2,2),
            nn.Flatten(),
            nn.Linear(64*4*4, 128),
            nn.ReLU(),
            nn.Linear(128,2),
            ])

    def forward(self, x):
        for layer in self.layers:
            x = layer(x)
        return x

# define CNN structure and its forward pass layer-by-layer
class AutoEncoder(nn.Module):
    def __init__(self):
        super().__init__()
        self.layers = nn.ModuleList([
            
            nn.Conv2d(in_channels=3,out_channels=16,kernel_size=3),
            nn.ReLU(),
            nn.Conv2d(in_channels=16,out_channels=32,kernel_size=3),
            nn.ReLU(),
            nn.Conv2d(in_channels=32,out_channels=64,kernel_size=3),
            nn.ReLU(),
            nn.Conv2d(in_channels=64,out_channels=128,kernel_size=3),
            nn.ReLU(),
            nn.ConvTranspose2d(in_channels=128,out_channels=64,kernel_size=3),
            nn.ReLU(),
            nn.ConvTranspose2d(in_channels=64,out_channels=32,kernel_size=3),
            nn.ReLU(),
            nn.ConvTranspose2d(in_channels=32,out_channels=16,kernel_size=3),
            nn.ReLU(),
            nn.ConvTranspose2d(in_channels=16,out_channels=3,kernel_size=3),
            nn.ReLU(),
            nn.Sigmoid()
            ]
        )

    def forward(self,x):
        for layer in self.layers:
            x = layer(x)
        return x

def train(model, device, train_loader, optimizer, criterion, epoch):
    model.train()

    for batch_idx, (data, target) in enumerate(train_loader):
        data, target = data.to(device), target.to(device)
        optimizer.zero_grad()
        output = model(data)
        loss = criterion(output, target)
        loss.backward()
        optimizer.step()
        if batch_idx % 10 == 0:
            print('Train Epoch: {} [{}/{} ({:.0f}%)]\tLoss: {:.6f}'.format(
                epoch, batch_idx * len(data), len(train_loader.dataset),
                100. * batch_idx / len(train_loader), loss.item()))


def test(model, device, test_loader, threshold):
    model.eval()
    correct = 0
    total = 0
    with torch.no_grad():
        for batch_idx, (data, target) in enumerate(test_loader):
            data,target  = data.to(device), target.to(device)
            output = model(data)
            
            pred = output.argmax(dim=1)

            correct += pred.eq(target).sum().item()
            total += target.size(0)

        print(f'Test acc: {100*correct/total:.2f}%')
            
def trainModel():
    device  = torch.device("cuda" if torch.cuda.is_available() else "cpu")
    transform = transforms.Compose([
        transforms.Resize((32,32)),
        transforms.ToTensor()
    ])


    train_dataset = IconDataset("../resources", transform=transform)
    train_loader = DataLoader(train_dataset, batch_size=64, shuffle=True)

    print(len(train_dataset))

    test_dataset = IconDataset("../resources", transform=transform)
    test_loader = DataLoader(test_dataset, batch_size=64, shuffle=False)

    model = TwoClass().to(device)
    #criterion = nn.MSELoss()
    criterion = nn.CrossEntropyLoss()
    optimizer = torch.optim.Adam(model.parameters(), lr=1e-3)

    for epoch in range(1,11):
        train(model, device, train_loader, optimizer, criterion, epoch)

    torch.save(model.state_dict(), "autoencoder.pth")

    '''
    model = AutoEncoder()
    model.load_state_dict(torch.load("autoencoder.pth"))
    model.to(device)
    ''' 

    test(model, device, test_loader, threshold=0.6)
    test(model, device, test_loader, threshold=0.7)
    test(model, device, test_loader, threshold=0.8)
 


class IconDataset(Dataset):
    def __init__(self, root_dir, transform=None):
        self.transform = transform
        self.samples = []

        skill_dir = os.path.join(root_dir, "skill-icons")
        for subfolder in os.listdir(skill_dir):
            subfolder_path = os.path.join(skill_dir, subfolder)
            if os.path.isdir(subfolder_path):
                for fname in os.listdir(subfolder_path):
                    if fname.lower().endswith((".png",".jpg",".jpeg")):
                        self.samples.append((os.path.join(subfolder_path,fname), 0))
        
        cooldown_dir = os.path.join(root_dir, "cooldown-icons")
        for fname in os.listdir(cooldown_dir):
            if fname.lower().endswith((".png",".jpg",".jpeg")):
                self.samples.append((os.path.join(cooldown_dir,fname), 1))

    def __len__(self):
        return len(self.samples)

    def __getitem__(self, idx):
        path, label = self.samples[idx]
        image = Image.open(path).convert("RGB")
        if self.transform:
            image = self.transform(image)
        return image, label


































if __name__ == "__main__":
    trainModel()                           




