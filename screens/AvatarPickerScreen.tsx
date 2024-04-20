import React, { useState } from 'react';
import { View, StyleSheet, TouchableOpacity, Image, ScrollView } from 'react-native';
import { Text, Button } from 'react-native-paper';
import { NativeStackScreenProps } from '@react-navigation/native-stack';
import { useDispatch } from 'react-redux';
import { setAvatarUri } from '../store';

type RootStackParamList = {
    SettingsScreen: {newAvatarUri: string};
    AvatarPickerScreen: undefined;
  };

type SettingsScreenProps = NativeStackScreenProps<RootStackParamList, 'AvatarPickerScreen'>;

type Avatar = {
    uri: string;
    image: NodeRequire;
}

const avatars: Avatar[] = [
    { uri: '../assets/avatars/avatar1.png', image: require('../assets/avatars/avatar1.png') },
    { uri: '../assets/avatars/avatar2.png', image: require('../assets/avatars/avatar2.png') },
    { uri: '../assets/avatars/avatar3.png', image: require('../assets/avatars/avatar3.png') },
    { uri: '../assets/avatars/avatar4.png', image: require('../assets/avatars/avatar4.png') },
    { uri: '../assets/avatars/avatar5.png', image: require('../assets/avatars/avatar5.png') },
    { uri: '../assets/avatars/avatar6.png', image: require('../assets/avatars/avatar6.png') },
    { uri: '../assets/avatars/avatar7.png', image: require('../assets/avatars/avatar7.png') },
    { uri: '../assets/avatars/avatar8.png', image: require('../assets/avatars/avatar8.png') },
    { uri: '../assets/avatars/avatar9.png', image: require('../assets/avatars/avatar9.png') },
];

let selectedAvatarUri: string = ''

const AvatarPickerScreen: React.FC<SettingsScreenProps> = ({ navigation }: SettingsScreenProps) => {
    const [selectedAvatar, setSelectedAvatar] = useState<string | null>(null);

    const handleAvatarSelect = (avatar: Avatar) => {
        setSelectedAvatar(avatar.uri);  // Store the URI of the selected avatar
    };

    const dispatch = useDispatch(); 

    const handleSave = (newAvatarUri: string) => {
        if (selectedAvatar) {
            let selectedAvatarUri: string = selectedAvatar // to prevent the "or null" state
            navigation.navigate('SettingsScreen', {newAvatarUri: selectedAvatarUri,}) //useRoute for local navigation
            dispatch(setAvatarUri(selectedAvatarUri)); //Redux for global storage 
        }
    };

    return (
        <ScrollView contentContainerStyle={styles.container}>
            <Text style={styles.title}>Choose Your Avatar</Text>
            <View style={styles.avatarContainer}>
                {avatars.map((avatar, index) => (
                    <TouchableOpacity key={index} onPress={() => handleAvatarSelect(avatar)}>
                        <Image
                            source={avatar.image}
                            style={[
                                styles.avatar,
                                selectedAvatar === avatar.uri ? styles.selectedAvatar : null,
                            ]}
                        />
                    </TouchableOpacity>
                ))}
            </View>
            <Button style={styles.button} mode="contained" onPress={() => handleSave(selectedAvatarUri)} disabled={!selectedAvatar}>
                Save
            </Button>
        </ScrollView>
    );
};

const styles = StyleSheet.create({
    container: {
        padding: 20,
        alignItems: 'center',
    },
    title: {
        fontSize: 20,
        marginBottom: 20,
    },
    button: {
        textAlign: 'center',
        alignSelf: 'center', 
      },
    avatarContainer: {
        flexDirection: 'row',
        flexWrap: 'wrap',
        justifyContent: 'center',
        marginBottom: 20,
    },
    avatar: {
        width: 70,
        height: 70,
        margin: 10,
        opacity: 0.6,
    },
    selectedAvatar: {
        opacity: 1,
        borderWidth: 3, 
        borderColor: '#4CAF50', //light green
        borderRadius: 35, 
    },
});

export default AvatarPickerScreen;